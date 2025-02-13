package main

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
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
		"pwd":  handlePwd,
		"cd":   handleCd,
	}

	scanner := bufio.NewScanner(os.Stdin)

	for {
		fmt.Fprint(os.Stdout, "$ ")

		if !scanner.Scan() {
			break
		}

		command := scanner.Text()
		words := parseInput(command)
		if len(words) == 0 {
			continue
		}
		cmd := words[0]
		arg := words[1:]
		if handler, exists := builtins[cmd]; exists {
			handler(arg)
		} else {
			runExternalCommand(cmd, arg)
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

func handlePwd(args []string) {
	path, err := os.Getwd()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error in getting path", err)
	} else {
		fmt.Println(path)
	}
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
			fmt.Printf("%s is %s\n", cmd, fullPath)
			found = true
			break
		}
	}
	if !found {
		fmt.Fprintf(os.Stderr, "%s: not found\n", cmd)
	}
}

func handleCd(args []string) {
	// Step 1: Determine the target directory
	var targetDir string
	if len(args) == 0 || args[0] == "~" {
		targetDir = os.Getenv("HOME")
	} else {
		targetDir = args[0]
	}
	// Step 2: convert the relative path to absolute path
	absPath, err := filepath.Abs(targetDir)
	if err != nil {
		fmt.Fprintf(os.Stderr, "cd: %v: No such file or directory\n", targetDir)
		return
	}

	// Step 3: checking if the path exists and its directory.
	info, err := os.Stat(absPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "cd: %s: No such file or directory\n", targetDir)
		return
	}
	if !info.IsDir() {
		fmt.Fprintf(os.Stderr, "cd: %s: No such file or directory\n", targetDir)
		return
	}

	// Step 4: Change the current working directory
	err = os.Chdir(absPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "cd: %s: %v\n", targetDir, err)
		return
	}
}

func parseInput(input string) []string {
	var args []string              // to store the final arguments
	var currentArg strings.Builder // to build each individual argument
	inSingleQuotes := false        // to track whether we are inside single quotes
	inDoubleQuotes := false        // to track whether we are inside Double quotes

	for i := 0; i < len(input); i++ {
		char := input[i]

		// handle the escape chars
		if char == '\\' && (inDoubleQuotes || inSingleQuotes) && i+1 < len(input) {
			i++ // move to the next char
			nextChar := input[i]

			// Only handle \, $, ", and newline inside double quotes
			if inDoubleQuotes && (nextChar == '\\' || nextChar == '$' || nextChar == '"' || nextChar == '\n') {
				currentArg.WriteByte(nextChar)
			} else {
				// outside double quotes or invalid escape sequence
				currentArg.WriteByte('\\')
				currentArg.WriteByte(nextChar)
			}
			continue
		}

		// handle single quotes
		if char == '\'' && !inDoubleQuotes {
			inSingleQuotes = !inSingleQuotes
			continue
		}
		// handle double quotes
		if char == '"' && !inSingleQuotes {
			inDoubleQuotes = !inDoubleQuotes
			continue
		}

		// handling space outside the quotes
		if !inSingleQuotes && !inDoubleQuotes && char == ' ' {
			if currentArg.Len() > 0 {
				args = append(args, currentArg.String())
				currentArg.Reset()
			}
			continue
		}
		currentArg.WriteByte(char)
	}

	// checking for unclosed quotes
	if inSingleQuotes {
		fmt.Fprintf(os.Stderr, "Error: unclosed single quote")
		return nil
	}
	if inDoubleQuotes {
		fmt.Fprintf(os.Stderr, "Error: unclosed double quote")
		return nil
	}

	if currentArg.Len() > 0 {
		args = append(args, currentArg.String())
	}
	return args
}

func fileExistsAndExecutable(path string) bool {
	info, err := os.Stat(path)
	if err != nil {
		return false
	}
	return !info.IsDir() && info.Mode().Perm()&0111 != 0
}

func runExternalCommand(cmd string, args []string) {
	// Step 1: Find the full path of the executable using PATH
	fullPath := findExecutableInPath(cmd)
	if fullPath == "" {
		fmt.Fprintf(os.Stderr, "%s: command not found\n", cmd)
		return
	}

	// Step 2: Use os/exec to run the program with its arguments
	output, err := executeProgram(fullPath, args)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error Executing %s: %v\n", cmd, err)
		return
	}

	// Print the output without adding an extra newline
	fmt.Print(string(output))

}

func findExecutableInPath(cmd string) string {
	// TODO: Implement logic to search for the executable in PATH
	// Use os.Getenv("PATH"), split it, and check each directory
	// Return the full path if found, otherwise return an empty string
	path := os.Getenv("PATH")
	if path == "" {
		return ""
	}
	dirs := strings.Split(path, ":")

	for _, dir := range dirs {
		fullPath := filepath.Join(dir, cmd)
		if fileExistsAndExecutable(fullPath) {
			return fullPath
		}
	}
	return ""
}

func executeProgram(program string, args []string) ([]byte, error) {
	// TODO: Use os/exec to execute the program with its arguments
	// Capture the combined output (stdout + stderr)
	// Return the output and any error that occurred
	basename := filepath.Base(program)
	cmd := exec.Command(program)
	cmd.Args = append([]string{basename}, args...)
	//cmd := exec.Command(program, append([]string{basename}, args...)...)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return nil, err
	}

	return output, nil
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
