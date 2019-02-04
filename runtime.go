package runtime

import (
	"bufio"
	"context"
	"github.com/gofunct/cflag"
	"github.com/gofunct/runtime/encoding"
	"io"
)

type Option func(r *Runtime)

// defaultIO is a basic implementation of the IO interface.
type Runtime struct {
	InputOutput *bufio.ReadWriter
	Scanner     *bufio.Scanner
	Encoders    encoding.EncoderGroup
	Decoders    encoding.DecoderGroup
	Handlers    []func(r *Runtime, ctx context.Context) error
	flagger 	*cflag.Flagger
}

func NewRuntime(opts ...Option) *Runtime {
	r := &Runtime{}
	for _, o := range opts {
		o(r)
	}
	return r
}

func NewDefaultRuntime(reader io.Reader, writer io.Writer) *Runtime {
	br := bufio.NewReader(reader)
	bw := bufio.NewWriter(writer)
	rw := bufio.NewReadWriter(br, bw)
	scan := bufio.NewScanner(rw.Reader)
	return &Runtime{InputOutput: rw, Scanner: scan, Encoders: encoding.DefaultEncoders, Decoders: encoding.DefaultDecoders}
}

func (r *Runtime) Close() error {
	return r.InputOutput.Flush()
}

func (r *Runtime) Read(p []byte) (n int, err error) {
	return r.InputOutput.Read(p)
}

func (r *Runtime) Write(p []byte) (n int, err error) {
	return r.InputOutput.Write(p)
}

func (r *Runtime) WriteTo(w io.Writer) (n int64, err error) {
	return r.InputOutput.WriteTo(w)
}

func (r *Runtime) ReadFrom(reader io.Reader) (n int64, err error) {
	return r.InputOutput.ReadFrom(reader)
}


func (r *Runtime) AddHandlerFunc(f func(r *Runtime, ctx context.Context) error) {
	r.Handlers = append(r.Handlers, f)
}

func (r *Runtime) Run(ctx context.Context) error {
	for _, f := range r.Handlers {
		if err := f(r, ctx); err != nil {
			return err
		}
	}
	return nil
}

func (r *Runtime) Runnable() bool {

	switch {
	case len(r.Handlers) > 0:
		return true

	default:
		return false
	}
}
