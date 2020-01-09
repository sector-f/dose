package main

import (
	"crypto/tls"
	"fmt"
	"net"
	"os"

	"github.com/spf13/pflag"
)

var (
	certFile    string
	keyFile     string
	certificate tls.Certificate

	tcpBind     []string
	tcpTlsBind  []string
	unixBind    []string
	unixTlsBind []string
)

func main() {
	pflag.StringVarP(&certFile, "cert", "", "", "TLS certificate file")
	pflag.StringVarP(&keyFile, "key", "", "", "TLS key file")
	pflag.StringArrayVarP(&tcpBind, "tcp", "", []string{}, "TCP socket to bind to")
	pflag.StringArrayVarP(&tcpTlsBind, "tcptls", "", []string{}, "TCP socket to bind to (with TLS)")
	pflag.StringArrayVarP(&unixBind, "unix", "", []string{}, "Unix socket to bind to")
	pflag.StringArrayVarP(&unixTlsBind, "unixtls", "", []string{}, "Unix socket to bind to (with TLS)")
	showHelp := pflag.BoolP("help", "h", false, "Show help message")
	pflag.Parse()

	if *showHelp {
		pflag.Usage()
		os.Exit(0)
	}

	if len(tcpBind) == 0 && len(tcpTlsBind) == 0 && len(unixBind) == 0 && len(unixTlsBind) == 0 {
		fmt.Fprintln(os.Stderr, "No TCP or Unix socket specified")
		os.Exit(1)
	}

	if (len(tcpTlsBind) > 0 || len(unixTlsBind) > 0) && (certFile == "" || keyFile == "") {
		fmt.Fprintln(os.Stderr, "Certificate and keyfile must be specified to use TLS")
		os.Exit(1)
	}

	if len(tcpTlsBind) > 0 || len(unixTlsBind) > 0 {
		cert, err := tls.LoadX509KeyPair(certFile, keyFile)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Certificate/key error: %s\n", err.Error())
			os.Exit(1)
		}
		certificate = cert
	}

	listeners := []*net.Listener{}

	for _, tcp := range tcpBind {
		listener, err := newTcpSocket(tcp, false)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
		listeners = append(listeners, listener)
	}

	for _, tcp := range tcpTlsBind {
		listener, err := newTcpSocket(tcp, true)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
		listeners = append(listeners, listener)
	}

	for _, unix := range unixBind {
		listener, err := newUnixSocket(unix, false)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
		listeners = append(listeners, listener)
	}

	for _, unix := range unixTlsBind {
		listener, err := newUnixSocket(unix, true)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
		listeners = append(listeners, listener)
	}

	runDownloadServer(listeners)
}
