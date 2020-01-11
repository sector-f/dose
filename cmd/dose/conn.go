package main

import (
	"crypto/tls"
	"errors"
	"net"
	"net/url"
	"syscall"
)

func tlsConfig() *tls.Config {
	return &tls.Config{
		Certificates: []tls.Certificate{certificate},
	}
}

func bind(u *url.URL) (*net.Listener, error) {
	switch u.Scheme {
	case "tcp":
		return newTcpSocket(u.Host, false)
	case "tcps":
		return newTcpSocket(u.Host, true)
	case "unix":
		return newUnixSocket(u.Path, false)
	case "unixs":
		return newUnixSocket(u.Path, true)
	default:
		return nil, errors.New("Unsupported connection type.\nMust be tcp[s]:// or unix[s]://")
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
