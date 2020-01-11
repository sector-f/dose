package main

import (
	"crypto/tls"
	"fmt"
	"net"
	"net/url"
	"os"

	"github.com/spf13/pflag"
)

var (
	certFile    string
	keyFile     string
	certificate tls.Certificate
	bindUrls    []string
)

func main() {
	pflag.StringVarP(&certFile, "cert", "", "", "TLS certificate file")
	pflag.StringVarP(&keyFile, "key", "", "", "TLS key file")
	pflag.StringArrayVarP(&bindUrls, "bind", "b", []string{}, "tcp[s]://addr:port or unix[s]:///path/to/socket")
	showHelp := pflag.BoolP("help", "h", false, "Show help message")
	pflag.CommandLine.SortFlags = false
	pflag.Parse()

	if *showHelp {
		pflag.Usage()
		os.Exit(0)
	}

	if len(bindUrls) == 0 {
		fmt.Fprintln(os.Stderr, "No TCP or Unix socket specified")
		os.Exit(1)
	}

	tlsCount := 0
	urls := []*url.URL{}
	for _, s := range bindUrls {
		url, err := url.Parse(s)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}

		switch url.Scheme {
		case "tcps", "unixs":
			tlsCount += 1
			urls = append(urls, url)
		case "tcp", "unix":
			urls = append(urls, url)
		default:
			fmt.Fprintln(os.Stderr, "Unsupported connection type.\nMust be tcp[s]:// or unix[s]://")
			os.Exit(1)
		}
	}

	if tlsCount > 0 {
		if certFile == "" || keyFile == "" {
			fmt.Fprintln(os.Stderr, "Certificate and keyfile must be specified to use TLS")
			os.Exit(1)
		}

		cert, err := tls.LoadX509KeyPair(certFile, keyFile)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Certificate/key error: %s\n", err.Error())
			os.Exit(1)
		}
		certificate = cert
	}

	listeners := []*net.Listener{}

	for _, url := range urls {
		listener, err := bind(url)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
		listeners = append(listeners, listener)
	}

	runDownloadServer(listeners)
}
