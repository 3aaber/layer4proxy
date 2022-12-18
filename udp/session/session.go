package session

import (
	"fmt"
	"net"
	"sync/atomic"
	"time"

	"layer4proxy/core"
)

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

func NewSession(clientAddr *net.UDPAddr, backEndConnection net.Conn, backend core.Upstream, cfg Config) *Session {

	s := &Session{
		cfg:        cfg,
		clientAddr: clientAddr,
		connection: backEndConnection,
		backend:    backend,
		outputC:    make(chan payload, MAX_PACKETS_QUEUE),
		stopC:      make(chan struct{}, 1),
	}

	s.closeChannelLoop()
	s.writeLoop()
	s.closeLoop()

	return s
}

func (s *Session) closeLoop() {
	go func() {

		var t *time.Timer

		if s.cfg.ClientIdleTimeout > 0 {
			t = time.NewTimer(s.cfg.ClientIdleTimeout)
		}

		for range t.C {
			s.Close()
		}

	}()

}

func (s *Session) writeLoop() {

	var t *time.Timer

	if s.cfg.ClientIdleTimeout > 0 {
		t = time.NewTimer(s.cfg.ClientIdleTimeout)
	}

	go func() {
		for pkt := range s.outputC {

			if t != nil {
				if !t.Stop() {
					<-t.C
				}
				t.Reset(s.cfg.ClientIdleTimeout)
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
		}
	}()
}

func (s *Session) closeChannelLoop() {

	var t *time.Timer

	if s.cfg.ClientIdleTimeout > 0 {
		t = time.NewTimer(s.cfg.ClientIdleTimeout)
	}
	go func() {
		for range s.stopC {

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

	}()
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
