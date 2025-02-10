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

		if command == "exit 0\n" {
			os.Exit(0)
		}
		word := strings.Fields(command)
		var words []string
		for _, w := range word {
			words = append(words, w)
		}
		if words[0] == "echo" {
			fmt.Println(strings.Join(words[1:], " "))
		} else {
			fmt.Println(command[:len(command)-1] + ": command not found")
		}
	}
}
