package throt

import (
	"context"
	"golang.org/x/time/rate"
	"io"
	"time"
)

type Writer struct {
	ctx      context.Context
	w        io.Writer
	limiters []*rate.Limiter
}

// NewWriter create io.Writer implementation with rate limitation.
func NewWriter(ctx context.Context, w io.Writer) *Writer {
	return &Writer{
		ctx: ctx,
		w:   w,
	}
}

// ApplyLimits set writing limit to bytePerSec bytes per second
func (th *Writer) ApplyLimits(limiters ...*Limiter) {
	for _, l := range limiters {
		l.AllowN(time.Now(), int(l.Limit())) // initialize a bucket with initial amount of tokens
		th.limiters = append(th.limiters, l.Limiter)
	}
}

// Write writes bytes from p.
func (th *Writer) Write(p []byte) (int, error) {
	if th.limiters == nil {
		return th.w.Write(p)
	}
	n, err := th.w.Write(p)
	if err != nil {
		return n, err
	}
	for _, l := range th.limiters {
		if err := l.WaitN(th.ctx, n); err != nil {
			return n, err
		}
	}
	return n, err
}
