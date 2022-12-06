package session

import (
	"sync"
	"time"
)

type Config struct {
	MaximumRequests    uint64
	MaximumResponses   uint64
	ClientIdleTimeout  time.Duration
	BackendIdleTimeout time.Duration
	Transparent        bool
}

const (
	UDP_PACKET_SIZE   = 65507
	MAX_PACKETS_QUEUE = 10000
)

var bufferPool = sync.Pool{
	New: func() interface{} {
		return make([]byte, UDP_PACKET_SIZE)
	},
}
