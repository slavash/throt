package throt

import (
	"context"
	"golang.org/x/time/rate"
	"io"
	"time"
)

// Reader io.Reader wrapper
type Reader struct {
	ctx      context.Context
	r        io.Reader
	limiters []*rate.Limiter
}

// NewReader create io.Reader implementation with rate limitation.
func NewReader(ctx context.Context, r io.Reader) *Reader {
	return &Reader{
		ctx: ctx,
		r:   r,
	}
}

// ApplyLimits set reading limit to bytePerSec bytes per second
func (th *Reader) ApplyLimits(limiters ...*rate.Limiter) {
	for _, l := range limiters {
		l.AllowN(time.Now(), int(l.Limit())) // initialize a bucket with initial amount of tokens
		th.limiters = append(th.limiters, l)
	}
}

// Read wrap reading with rate limitation
func (th *Reader) Read(p []byte) (int, error) {
	if th.limiters == nil || len(th.limiters) == 0 {
		return th.r.Read(p)
	}
	n, err := th.r.Read(p)
	if err != nil {
		return n, err
	}
	for _, l := range th.limiters {
		if err := l.WaitN(th.ctx, n); err != nil {
			return n, err
		}
	}

	return n, nil
}
