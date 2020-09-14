package throt

import (
	"context"
	"errors"
	"github.com/stretchr/testify/assert"
	"golang.org/x/time/rate"
	"io"
	"reflect"
	"testing"
	"time"
)

type ReaderMock struct {
	Payload func([]byte) int
	Error   error
}

func (r *ReaderMock) Read(b []byte) (int, error) {
	return r.Payload(b), r.Error
}

func TestNewReader(t *testing.T) {
	type args struct {
		ctx context.Context
		r   io.Reader
	}
	tests := []struct {
		name string
		args args
		want *Reader
	}{
		{
			"TestNewReader",
			args{
				ctx: context.Background(),
				r:   &ReaderMock{},
			},
			&Reader{
				ctx:     context.Background(),
				r:       &ReaderMock{},
				limiter: nil,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewReader(tt.args.ctx, tt.args.r); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewReader() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestReader_ApplyLimit(t *testing.T) {
	type fields struct {
		ctx     context.Context
		r       io.Reader
		limiter *rate.Limiter
	}
	type args struct {
		l *Limiter
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		{
			"TestApplyLimit",
			fields{
				ctx:     context.Background(),
				r:       &ReaderMock{},
				limiter: rate.NewLimiter(22, 2),
			},
			args{l: NewLimiter(10, 1)},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			th := &Reader{
				ctx:     tt.fields.ctx,
				r:       tt.fields.r,
				limiter: tt.fields.limiter,
			}
			assert.Equal(t, 22, int(th.limiter.Limit()))
			th.ApplyLimit(tt.args.l)
			assert.Equal(t, 10, int(th.limiter.Limit()))
		})
	}
}

func TestReader_Read(t *testing.T) {
	type fields struct {
		ctx     context.Context
		r       io.Reader
		limiter *rate.Limiter
	}
	type args struct {
		p []byte
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		delay   time.Duration
		want    int
		wantErr bool
	}{
		{
			"TestRead buffer lt burst",
			fields{
				ctx: context.Background(),
				r: &ReaderMock{
					Payload: func(bytes []byte) int {
						return len(bytes)
					},
					Error: nil,
				},
				limiter: rate.NewLimiter(100, 10),
			},
			args{
				p: []byte("data"),
			},
			0,
			4, false,
		},
		{
			"TestRead  buffer gt burst",
			fields{
				ctx: context.Background(),
				r: &ReaderMock{
					Payload: func(bytes []byte) int {
						return len(bytes)
					},
					Error: nil,
				},
				limiter: rate.NewLimiter(100, 10),
			},
			args{
				p: []byte("datadatadata"),
			},
			0,
			12, true,
		},
		{
			"TestRead  buffer eq burst",
			fields{
				ctx: context.Background(),
				r: &ReaderMock{
					Payload: func(bytes []byte) int {
						return len(bytes)
					},
					Error: nil,
				},
				limiter: rate.NewLimiter(100, 12),
			},
			args{
				p: []byte("datadatadata"),
			},
			0,
			12, false,
		},
		{
			"TestRead  buffer eq burst",
			fields{
				ctx: context.Background(),
				r: &ReaderMock{
					Payload: func(bytes []byte) int {
						return len(bytes)
					},
					Error: nil,
				},
				limiter: nil,
			},
			args{
				p: []byte("datadatadata"),
			},
			0,
			12, false,
		},
		{
			"TestRead with error",
			fields{
				ctx: context.Background(),
				r: &ReaderMock{
					Payload: func(bytes []byte) int {
						return len(bytes)
					},
					Error: errors.New("Some error"),
				},
				limiter: rate.NewLimiter(100, 10),
			},
			args{
				p: []byte("datadatadata"),
			},
			0,
			12, true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			th := &Reader{
				ctx:     tt.fields.ctx,
				r:       tt.fields.r,
				limiter: tt.fields.limiter,
			}
			got, err := th.Read(tt.args.p)
			if (err != nil) != tt.wantErr {
				t.Errorf("Read() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("Read() got = %v, want %v", got, tt.want)
			}
		})
	}
}
