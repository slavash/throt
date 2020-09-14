package throt

import (
	"context"
	"golang.org/x/time/rate"
	"io"
)

// Reader io.Reader wrapper
type Reader struct {
	ctx     context.Context
	r       io.Reader
	limiter *rate.Limiter
}

// NewReader create io.Reader implementation with rate limitation.
func NewReader(ctx context.Context, r io.Reader) *Reader {
	return &Reader{
		ctx: ctx,
		r:   r,
	}
}

// ApplyLimit set reading limit to bytePerSec bytes per second
func (th *Reader) ApplyLimit(l *Limiter) {
	th.limiter = l.Limiter
}

// Read wrap reading with rate limitation
func (th *Reader) Read(p []byte) (int, error) {
	if th.limiter == nil {
		return th.r.Read(p)
	}
	n, err := th.r.Read(p)
	if err != nil {
		return n, err
	}
	if err := th.limiter.WaitN(th.ctx, n); err != nil {
		return n, err
	}

	return n, nil
}
