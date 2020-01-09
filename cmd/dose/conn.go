package main

import (
	"crypto/tls"
	"net"
	"syscall"
)

func tlsConfig() *tls.Config {
	return &tls.Config{
		Certificates: []tls.Certificate{certificate},
	}
}

// Creates new Unix socket with 0600 permissions
func newUnixSocket(path string, useTLS bool) (*net.Listener, error) {
	var (
		l   net.Listener
		err error
	)

	oldUmask := syscall.Umask(0177)
	defer syscall.Umask(oldUmask)

	if useTLS {
		l, err := tls.Listen("unix", path, tlsConfig())
		if err != nil {
			return nil, err
		}

		return &l, nil
	}

	l, err = net.Listen("unix", path)
	if err != nil {
		return nil, err
	}

	return &l, nil
}

func newTcpSocket(addr string, useTLS bool) (*net.Listener, error) {
	var (
		l   net.Listener
		err error
	)

	oldUmask := syscall.Umask(0177)
	defer syscall.Umask(oldUmask)

	if useTLS {
		l, err := tls.Listen("tcp", addr, tlsConfig())
		if err != nil {
			return nil, err
		}

		return &l, nil
	}

	l, err = net.Listen("tcp", addr)
	if err != nil {
		return nil, err
	}

	return &l, nil
}
