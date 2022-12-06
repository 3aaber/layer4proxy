package session

import "sync"

type payload struct {
	buffer []byte
	length int
}

var bufferPool = sync.Pool{
	New: func() interface{} {
		return make([]byte, UDP_PACKET_SIZE)
	},
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
