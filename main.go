package main

import (
	"fmt"
	"io"
	"os"
)

const PROGRAM_NAME = "gtsort"
const VERSION = "1"

func usage(status int) {
	if status != 0 {
		fmt.Fprintf(os.Stderr, "Try '%s --help' for more information.\n", PROGRAM_NAME)
	} else {
		fmt.Printf("Usage: %s [OPTION] [FILE]\n", PROGRAM_NAME)
		fmt.Println("Write totally ordered list consistent with the partial ordering in FILE.")
		fmt.Println("With no FILE, or when FILE is -, read standard input.")
		fmt.Println("  --help     display this help and exit")
		fmt.Println("  --version  output version information and exit")
	}
	os.Exit(status)
}

func main() {
	args := os.Args[1:]
	if len(args) > 0 {
		switch args[0] {
		case "--help":
			usage(0)
		case "--version":
			fmt.Printf("%s version %s\n", PROGRAM_NAME, VERSION)
			os.Exit(0)
		}
	}

	if len(args) > 1 {
		fmt.Fprintf(os.Stderr, "extra operand %s\n", args[1])
		usage(1)
	}

	var input io.Reader = os.Stdin
	if len(args) == 1 && args[0] != "-" {
		f, err := os.Open(args[0])
		if err != nil {
			fmt.Fprintf(os.Stderr, "%s: %v\n", args[0], err)
			os.Exit(1)
		}
		defer f.Close()
		input = f
	}

	if err := tsort(input); err != nil {
		fmt.Fprintf(os.Stderr, "%s: %v\n", PROGRAM_NAME, err)
		os.Exit(1)
	}
}
