package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

// Ensures gofmt doesn't remove the "fmt" import in stage 1 (feel free to remove this!)
var _ = fmt.Fprint

func main() {
	for {
		fmt.Fprint(os.Stdout, "$ ")

		// Wait for user input
		command, err := bufio.NewReader(os.Stdin).ReadString('\n')
		if err != nil {
			fmt.Fprintln(os.Stderr, "Error reading input:", err)
			os.Exit(1)
		}
		word := strings.Fields(command)
		var words []string
		for _, w := range word {
			words = append(words, w)
		}

		switch words[0] {
		case "echo":
			fmt.Fprintf(os.Stdout, "%s\n", strings.Join(words[1:], " "))
		case "exit":
			os.Exit(0)
		default:
			fmt.Fprintf(os.Stderr, "%s: command not found\n", command[:len(command)-1])
		}

	}
}
