package clientgrpc

import (
	"fmt"

	"google.golang.org/grpc"
)

// NewClientConn creates a client connection to the given addr.
func NewClientConn(
	addr string,
	opt ...Option,
) (
	_ *grpc.ClientConn,
	_ error,
) {
	opts := defaultOptions()

	for _, o := range opt {
		o.apply(opts)
	}

	conn, err := grpc.Dial(addr, opts.dialOptions...)
	if err != nil {
		return nil, fmt.Errorf("grpc dial: %w", err)
	}

	return conn, nil
}
