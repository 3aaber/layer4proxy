package session

import (
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
