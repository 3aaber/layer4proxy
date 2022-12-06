package session

import "time"

type Config struct {
	MaximumRequests    uint64
	MaximumResponses   uint64
	ClientIdleTimeout  time.Duration
	BackendIdleTimeout time.Duration
	Transparent        bool
}
