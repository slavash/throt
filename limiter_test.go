package throt

import (
	"reflect"
	"testing"

	"golang.org/x/time/rate"
)

func TestNewLimiter(t *testing.T) {
	type args struct {
		rateLimit int
		burst     int
	}
	tests := []struct {
		name string
		args args
		want *Limiter
	}{
		{
			"TestNewLimiter",
			args{
				rateLimit: 100,
				burst:     1,
			},
			&Limiter{
				rate.NewLimiter(100, 1),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewLimiter(tt.args.rateLimit, tt.args.burst); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewLimiter() = %v, want %v", got, tt.want)
			}
		})
	}
}
