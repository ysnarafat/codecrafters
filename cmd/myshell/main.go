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
var validCommands = []string{"echo", "exit", "type", "pwd", "cd", "cat"}

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
			content := strings.TrimSpace(strings.TrimPrefix(input, "echo "))

			if strings.HasPrefix(content, "'") && strings.HasSuffix(content, "'") {
				content = content[1 : len(content)-1]
			} else {
				words := strings.Fields(content) // This will handle multiple spaces
				content = strings.Join(words, " ")
			}

			fmt.Println(content)
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

		if strings.HasPrefix(input, "cd ") {
			changeToDirectory := strings.TrimPrefix(input, "cd ")

			if changeToDirectory == "~" {
				changeToDirectory, _ = os.UserHomeDir()
			}

			err := os.Chdir(changeToDirectory)

			if err != nil {
				fmt.Printf("cd: %s: No such file or directory\n", changeToDirectory)
			}

			continue
		}

		if strings.HasPrefix(input, "cat ") {
			concatArgs := strings.TrimPrefix(input, "cat ")
			args := strings.Split(concatArgs, " ")

			for _, item := range args {
				item = strings.Trim(item, "'")
				content, err := os.ReadFile(item)

				if err == nil {
					fmt.Printf("%s ", string(content))
				}
			}

			fmt.Println()

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
	for _, item := range validCommands {
		if item == str {
			return true
		}
	}
	return false
}
