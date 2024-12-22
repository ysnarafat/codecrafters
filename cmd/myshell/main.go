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

		message := input + ": command not found"

		fmt.Println(message)
	}
}
