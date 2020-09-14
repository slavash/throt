package throt

import (
	"bytes"
	"context"
	"golang.org/x/time/rate"
	"io"
	"reflect"
	"testing"
)

// TODO Add tests...
func TestNewWriter(t *testing.T) {
	type args struct {
		ctx context.Context
	}
	tests := []struct {
		name  string
		args  args
		wantW string
		want  *Writer
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := &bytes.Buffer{}
			got := NewWriter(tt.args.ctx, w)
			if gotW := w.String(); gotW != tt.wantW {
				t.Errorf("NewWriter() gotW = %v, want %v", gotW, tt.wantW)
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewWriter() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestWriter_ApplyLimit(t *testing.T) {
	type fields struct {
		ctx     context.Context
		w       io.Writer
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
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_ = &Writer{
				ctx:     tt.fields.ctx,
				w:       tt.fields.w,
				limiter: tt.fields.limiter,
			}
		})
	}
}

func TestWriter_Write(t *testing.T) {
	type fields struct {
		ctx     context.Context
		w       io.Writer
		limiter *rate.Limiter
	}
	type args struct {
		p []byte
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    int
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			th := &Writer{
				ctx:     tt.fields.ctx,
				w:       tt.fields.w,
				limiter: tt.fields.limiter,
			}
			got, err := th.Write(tt.args.p)
			if (err != nil) != tt.wantErr {
				t.Errorf("Write() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("Write() got = %v, want %v", got, tt.want)
			}
		})
	}
}
