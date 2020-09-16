package throt

import (
	"context"
	"io"

	"golang.org/x/time/rate"
)

type Writer struct {
	ctx     context.Context
	w       io.Writer
	limiter *rate.Limiter
}

// NewWriter create io.Writer implementation with rate limitation.
func NewWriter(ctx context.Context, w io.Writer) *Writer {
	return &Writer{
		ctx: ctx,
		w:   w,
	}
}

// ApplyLimit set writing limit to bytePerSec bytes per second
func (th *Writer) ApplyLimit(l *Limiter) {
	th.limiter = l.Limiter
}

// Write writes bytes from p.
func (th *Writer) Write(p []byte) (int, error) {
	if th.limiter == nil {
		return th.w.Write(p)
	}
	n, err := th.w.Write(p)
	if err != nil {
		return n, err
	}

	if err := th.limiter.WaitN(th.ctx, n); err != nil {
		return n, err
	}

	return n, err
}
