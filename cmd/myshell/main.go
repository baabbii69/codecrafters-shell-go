package main

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

var builtins map[string]func([]string)

// Ensures gofmt doesn't remove the "fmt" import in stage 1 (feel free to remove this!)
var _ = fmt.Fprint

func main() {
	builtins = map[string]func([]string){
		"echo": handleEcho,
		"exit": handleExit,
		"type": handleType,
	}

	scanner := bufio.NewScanner(os.Stdin)

	for {
		fmt.Fprint(os.Stdout, "$ ")

		if !scanner.Scan() {
			break
		}

		command := scanner.Text()
		words := strings.Fields(command)
		if len(words) == 0 {
			continue
		}
		cmd := words[0]
		arg := words[1:]
		if handler, exists := builtins[cmd]; exists {
			handler(arg)
		} else {
			fmt.Fprintf(os.Stderr, "%s: command not found\n", cmd)
		}

		if err := scanner.Err(); err != nil {
			fmt.Fprintf(os.Stderr, "Error reading input:", err)
			os.Exit(1)
		}
	}
}
func handleEcho(args []string) {
	fmt.Println(strings.Join(args, " "))
}

func handleExit(args []string) {
	os.Exit(0)
}

func handleType(args []string) {
	if len(args) == 0 {
		fmt.Fprintln(os.Stderr, "type: missing argument")
		return
	}

	cmd := args[0]
	//checks if the command is builtin or not by checking it from builtin slice
	if _, exists := builtins[cmd]; exists {
		fmt.Printf("%s is a shell builtin\n", cmd)
		return
	}
	//get the path environment variable
	path := os.Getenv("PATH")
	if path == "" {
		fmt.Fprintf(os.Stderr, "%s: PATH not set\n", cmd)
		return
	}
	//split the path into directories
	dirs := strings.Split(path, ":")
	found := false
	for _, dir := range dirs {
		fullPath := filepath.Join(dir, cmd)
		if fileExistsAndExecutable(fullPath) {
			fmt.Printf("%s is in %s\n", cmd, fullPath)
			found = true
			break
		}
	}
	if !found {
		fmt.Fprintf(os.Stderr, "%s: not found\n", cmd)
	}
}

func fileExistsAndExecutable(path string) bool {
	info, err := os.Stat(path)
	if err != nil {
		return false
	}
	return !info.IsDir() && info.Mode().Perm()&0111 != 0
}

//fmt.Fprint(os.Stdout, "$ ")
//
//// Wait for user input
//command, err := bufio.NewReader(os.Stdin).ReadString('\n')
//if err != nil {
//	fmt.Fprintln(os.Stderr, "Error reading input:", err)
//	os.Exit(1)
//}
//word := strings.Fields(command)
//var words []string
//for _, w := range word {
//	words = append(words, w)
//}
//
//switch words[0] {
//case "echo":
//	fmt.Fprintf(os.Stdout, "%s\n", strings.Join(words[1:], " "))
//case "exit":
//	os.Exit(0)
//default:
//	fmt.Fprintf(os.Stderr, "%s: command not found\n", command[:len(command)-1])
//}
