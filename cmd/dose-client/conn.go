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
		InsecureSkipVerify: allowInsecure,
		// Certificates: []tls.Certificate{certificate},
	}
}

func bind(u *url.URL) (net.Conn, error) {
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
func newUnixSocket(path string, useTLS bool) (net.Conn, error) {
	var (
		c   net.Conn
		err error
	)

	oldUmask := syscall.Umask(0177)
	defer syscall.Umask(oldUmask)

	if useTLS {
		c, err := tls.Dial("unix", path, tlsConfig())
		if err != nil {
			return nil, err
		}

		return c, nil
	}

	c, err = net.Dial("unix", path)
	if err != nil {
		return nil, err
	}

	return c, nil
}

func newTcpSocket(addr string, useTLS bool) (net.Conn, error) {
	var (
		c   net.Conn
		err error
	)

	oldUmask := syscall.Umask(0177)
	defer syscall.Umask(oldUmask)

	if useTLS {
		c, err := tls.Dial("tcp", addr, tlsConfig())
		if err != nil {
			return nil, err
		}

		return c, nil
	}

	c, err = net.Dial("tcp", addr)
	if err != nil {
		return nil, err
	}

	return c, nil
}
