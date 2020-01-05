package main

import (
	"fmt"
	"os"
)

func printHelp() {
	fmt.Printf("Usage:\n")
	fmt.Printf("%v add URL PATH\n", os.Args[0])
	fmt.Printf("%v cancel PATH\n", os.Args[0])
}

func main() {
	args := os.Args

	if len(args) < 2 {
		printHelp()
		os.Exit(1)
	}

	switch args[1] {
	case "add":
		if len(args) != 4 {
			printHelp()
			os.Exit(1)
		}

		url := args[2]
		filepath := args[3]

		if url == "" || filepath == "" {
			printHelp()
			os.Exit(1)
		}

		download(url, filepath)
	case "cancel":
		if len(args) != 3 {
			printHelp()
			os.Exit(1)
		}

		filepath := args[2]

		if filepath == "" {
			printHelp()
			os.Exit(1)
		}

		cancel(filepath)
	case "help", "-h", "--help":
		printHelp()
		os.Exit(0)
	}
}
