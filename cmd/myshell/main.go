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
				fmt.Println(commandType + ": not found")
			}

			continue
		}

		handleInvalidCommand(input)
	}
}

func handleInvalidCommand(input string) {
	message := input + ": command not found"
	fmt.Println(message)
}

func isValidCommand(str string) bool {
	validCommands := []string{"echo", "exit"}

    for _, item := range validCommands {
        if item == str {
            return true
        }
    }
    return false
}

