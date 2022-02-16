package net

import (
	"errors"
	"fmt"
	"net"
)

// ErrInvalidAddrType is returned if an invalid listener address type
// is used.
var ErrInvalidAddrType = errors.New("invalid address type")

// GetFreePort returns an available open port on the current machine.
// Mostly used for testing.
func GetFreePort() (int, error) {
	addr, err := net.ResolveTCPAddr("tcp", "localhost:0")
	if err != nil {
		return 0, fmt.Errorf("resolve tcp addr: %w", err)
	}

	l, err := net.ListenTCP("tcp", addr)
	if err != nil {
		return 0, fmt.Errorf("listen tcp: %w", err)
	}

	err = l.Close()
	if err != nil {
		return 0, fmt.Errorf("close listener: %w", err)
	}

	addr, ok := l.Addr().(*net.TCPAddr)
	if !ok {
		return 0, ErrInvalidAddrType
	}

	return addr.Port, nil
}
