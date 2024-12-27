package main

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

// Ensures gofmt doesn't remove the "fmt" import in stage 1 (feel free to remove this!)
var _ = fmt.Fprint

func main() {
	for {
		fmt.Fprint(os.Stdout, "$ ")

		// Wait for user input
		input, error := bufio.NewReader(os.Stdin).ReadString('\n')

		if error != nil {
			fmt.Println(error)
			return
		}

		input = strings.TrimSpace(input)

		if input == "exit 0" {
			break
		}

		if strings.HasPrefix(input, "echo ") {
			echoMessage := strings.TrimSpace(strings.TrimPrefix(input, "echo "))
			fmt.Println(echoMessage)
			continue
		}

		if strings.HasPrefix(input, "type ") {
			commandType := strings.TrimSpace(strings.TrimPrefix(input, "type "))

			if isValidCommand(commandType) {
				fmt.Println(commandType + " is a shell builtin")
			} else {
				handleBuiltInCommand(commandType)
			}

			continue
		}

		if input == "pwd" {
			pwd, err := os.Getwd()
			if err == nil {
				fmt.Println(pwd)
			}

			continue
		}

		handleInvalidCommand(input)
	}
}

func handleBuiltInCommand(commandType string) {
	paths := strings.Split(os.Getenv("PATH"), ":")

	for _, path := range paths {
		fp := filepath.Join(path, commandType)
		if _, err := os.Stat(fp); err == nil {
			fmt.Printf("%s is %s\n", commandType, fp)
			return
		}
	}

	fmt.Printf("%s: not found\n", commandType)
}

func handleInvalidCommand(input string) {
	cmds := strings.Split(input, " ")
	command := exec.Command(cmds[0], cmds[1:]...)

	 command.Stderr = os.Stderr
	 command.Stdout = os.Stdout

	err := command.Run()
	if err != nil {
		fmt.Printf("%s: command not found\n", cmds[0])
	}
}

func isValidCommand(str string) bool {
	validCommands := []string{"echo", "exit", "type"}

	for _, item := range validCommands {
		if item == str {
			return true
		}
	}
	return false
}
