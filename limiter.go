package throt

import (
	"golang.org/x/time/rate"
)

// Limiter rate.Limiter wrapper
type Limiter struct {
	*rate.Limiter
}

// NewLimiter create new instance of Limiter with the rateLimit bytes per second limit rate
func NewLimiter(rateLimit int, burst int) *Limiter {
	return &Limiter{rate.NewLimiter(rate.Limit(rateLimit), burst)}
}
