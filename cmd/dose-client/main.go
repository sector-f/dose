package main

import (
	"bufio"
	"fmt"
	"net/url"
	"os"
	"syscall"

	"github.com/sector-f/dose"
	"github.com/spf13/pflag"
	"golang.org/x/crypto/ssh/terminal"
)

var (
	socketAddr    string
	allowInsecure bool
	noAuth        bool
)

func printHelp() {
	fmt.Printf("Usage:\n")
	fmt.Printf("%v add URL PATH\n", os.Args[0])
	fmt.Printf("%v cancel PATH\n", os.Args[0])
}

func main() {
	pflag.StringVarP(&socketAddr, "bind", "b", "unix:///tmp/dose.socket", "tcp[s]:// or unix[s]:// addr to connect to")
	pflag.BoolVarP(&allowInsecure, "insecure", "k", false, "Accept certificates from server even if they can't be validated")
	pflag.BoolVarP(&noAuth, "noauth", "n", false, "Skip sending a username/password to the server")
	showHelp := pflag.BoolP("help", "h", false, "Show help message")
	pflag.Parse()

	if *showHelp {
		pflag.Usage()
		os.Exit(0)
	}

	url, err := url.Parse(socketAddr)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	switch url.Scheme {
	case "tcp", "tcps", "unix", "unixs":
	default:
		fmt.Fprintln(os.Stderr, "Unsupported connection type.\nMust be tcp[s]:// or unix[s]://")
		os.Exit(1)
	}

	args := pflag.Args()
	if len(args) < 1 {
		printHelp()
		os.Exit(1)
	}

	conn, err := bind(url)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	defer conn.Close()

	if !noAuth {
		reader := bufio.NewReader(os.Stdin)
		fmt.Print("Enter username: ")
		username, _ := reader.ReadString('\n')
		fmt.Print("Enter password: ")
		password, _ := terminal.ReadPassword(int(syscall.Stdin))
		fmt.Println()
		dose.WriteMessage(conn, dose.AuthRequest{username, string(password)})
	}

	switch args[0] {
	case "add":
		if len(args) != 3 {
			printHelp()
			os.Exit(1)
		}

		url := args[1]
		filepath := args[2]

		if url == "" || filepath == "" {
			printHelp()
			os.Exit(1)
		}

		download(conn, url, filepath)
	case "cancel":
		if len(args) != 2 {
			printHelp()
			os.Exit(1)
		}

		filepath := args[1]

		if filepath == "" {
			printHelp()
			os.Exit(1)
		}

		cancel(conn, filepath)
	case "serverinfo":
		getServerInfo(conn)
	case "help", "-h", "--help":
		printHelp()
		os.Exit(0)
	}
}
