package clientgrpc

import (
	"context"
	"net"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type options struct {
	dialOptions []grpc.DialOption
}

func defaultOptions() *options {
	return &options{
		dialOptions: []grpc.DialOption{},
	}
}

// Option configures how we set up the connection.
type Option interface {
	apply(*options)
}

type funcServerOption struct {
	f func(*options)
}

func (f *funcServerOption) apply(do *options) {
	f.f(do)
}

func newFuncServerOption(f func(*options)) *funcServerOption {
	return &funcServerOption{
		f: f,
	}
}

// WithNoTLS disables transport security for the client.
// Replacement for grpc.WithInsecure().
func WithNoTLS() Option {
	return newFuncServerOption(func(o *options) {
		o.dialOptions = append(
			o.dialOptions,
			grpc.WithTransportCredentials(insecure.NewCredentials()),
		)
	})
}

// WithContextDialer wraps the grpc.WithContextDialer option.
func WithContextDialer(
	d func(context.Context, string) (net.Conn, error),
) Option {
	return newFuncServerOption(func(o *options) {
		o.dialOptions = append(
			o.dialOptions,
			grpc.WithContextDialer(d),
		)
	})
}
