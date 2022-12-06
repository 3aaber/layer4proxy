package session

import (
	"fmt"
	"net"
	"sync"
	"sync/atomic"
	"time"

	"layer4proxy/core"
)

const (
	UDP_PACKET_SIZE   = 65507
	MAX_PACKETS_QUEUE = 10000
)

var bufferPool = sync.Pool{
	New: func() interface{} {
		return make([]byte, UDP_PACKET_SIZE)
	},
}

type payload struct {
	buffer []byte
	length int
}

func (p payload) buf() []byte {
	if p.buffer == nil {
		return nil
	}

	return p.buffer[0:p.length]
}

func (p payload) release() {
	if p.buffer == nil {
		return
	}
	bufferPool.Put(p.buffer)
}

type Session struct {
	sent uint64
	recv uint64

	cfg Config

	clientAddr *net.UDPAddr
	connection net.Conn
	backend    core.Upstream

	outputC chan payload
	stopC   chan struct{}
	stopped uint32
}

func NewSession(clientAddr *net.UDPAddr, conn net.Conn, backend core.Upstream, cfg Config) *Session {

	s := &Session{
		cfg:        cfg,
		clientAddr: clientAddr,
		connection: conn,
		backend:    backend,
		outputC:    make(chan payload, MAX_PACKETS_QUEUE),
		stopC:      make(chan struct{}, 1),
	}

	go func() {

		var t *time.Timer
		var tC <-chan time.Time

		if cfg.ClientTimeout > 0 {
			t = time.NewTimer(cfg.ClientTimeout)
			tC = t.C
		}

		for {
			select {

			case <-tC:
				s.Close()
			case pkt := <-s.outputC:
				if t != nil {
					if !t.Stop() {
						<-t.C
					}
					t.Reset(cfg.ClientTimeout)
				}

				if pkt.buffer == nil {
					panic("Program error, output channel should not be closed here")
				}

				n, err := s.connection.Write(pkt.buf())
				pkt.release()

				if err != nil {
					fmt.Printf("Could not write data to udp connection: %v", err)
					break
				}

				if n != pkt.length {
					fmt.Printf("Short write error: should write %d bytes, but %d written", pkt.length, n)
					break
				}

				if s.cfg.MaximumRequests > 0 && atomic.AddUint64(&s.sent, 1) > s.cfg.MaximumRequests {
					fmt.Printf("Restricted to send more UDP packets")
					break
				}
			case <-s.stopC:
				atomic.StoreUint32(&s.stopped, 1)
				if t != nil {
					t.Stop()
				}
				s.connection.Close()

				// drain output packets channel and free buffers
				for {
					select {
					case pkt := <-s.outputC:
						pkt.release()
					default:
						return
					}
				}

			}
		}

	}()

	return s
}

func (s *Session) Write(buf []byte) error {
	if atomic.LoadUint32(&s.stopped) == 1 {
		return fmt.Errorf("closed session")
	}

	dup := bufferPool.Get().([]byte)
	n := copy(dup, buf)

	select {
	case s.outputC <- payload{dup, n}:
	default:
		bufferPool.Put(dup)
	}

	return nil
}

/**
 * ListenResponses waits for responses from backend, and sends them back to client address via
 * server connection, so that client is not confused with source host:port of the
 * packet it receives
 */
func (s *Session) ListenResponses(sendTo *net.UDPConn) {

	go func() {
		b := make([]byte, UDP_PACKET_SIZE)

		defer s.Close()

		for {

			if s.cfg.BackendIdleTimeout > 0 {
				s.connection.SetReadDeadline(time.Now().Add(s.cfg.BackendIdleTimeout))
			}

			n, err := s.connection.Read(b)

			if err != nil {
				if err, ok := err.(net.Error); ok && err.Timeout() {
					return
				}

				if atomic.LoadUint32(&s.stopped) == 0 {
					fmt.Printf("Failed to read from backend: %v", err)
				}
				return
			}

			m, err := sendTo.WriteToUDP(b[0:n], s.clientAddr)

			if err != nil {
				fmt.Printf("Could not send backend response to client: %v", err)
				return
			}

			if m != n {
				return
			}

			if s.cfg.MaximumResponses > 0 && atomic.AddUint64(&s.recv, 1) >= s.cfg.MaximumResponses {
				return
			}
		}
	}()
}

func (s *Session) IsDone() bool {
	return atomic.LoadUint32(&s.stopped) == 1
}

func (s *Session) Close() {
	select {
	case s.stopC <- struct{}{}:
	default:
	}
}