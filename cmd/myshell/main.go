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
var validCommands = []string{"echo", "exit", "type", "pwd", "cd"}

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
			arg := strings.TrimPrefix(input, "echo ")
			result := handleEchoCommand(arg)
			//content := strings.Join(result, " ")
			fmt.Println(result)
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

			var args []string

			if strings.HasPrefix(concatArgs, "'") {
				args = parseSingleQuotedStrings(concatArgs)
			} else {
				args = parseDoubleQuotedStrings(concatArgs)
			}

			for _, item := range args {
				//items := strings.Split(item, "'")
				// var path string = ""
				// for _, ch := range item {
				// 	strings.Re
				// }

				path := item
				//path := strings.ReplaceAll(item, "'", "")

				//path := handleEchoCommand(item)
				//fmt.Println(path)

				content, err := os.ReadFile(path)

				if err == nil {
					fmt.Printf("%s", string(content))
				}
			}

			continue
		}

		handleInvalidCommand(input)
	}
}

func handleEchoCommand(input string) string {
	printAble := input
	splitter := " "
	printAbleTemp := ""
	inSingleQuote := false
	inDoubleQuote := false

	if strings.HasPrefix(printAble, "'") && strings.HasSuffix(printAble, "'") {
		splitter = "'"
		inSingleQuote = true
	} else if strings.HasPrefix(printAble, "\"") && strings.HasSuffix(printAble, "\"") {
		splitter = "\""
		inDoubleQuote = true
	}

	if(!(inSingleQuote || inDoubleQuote)) {
		if(strings.Contains(printAble, "\\")) {
			return strings.ReplaceAll(printAble, "\\", "")
		}

		items := strings.Split(printAble, " ")

		filteredItems := Filter(items, func(s string) bool {
			return s != ""
		})

		return strings.TrimSpace(strings.Join(filteredItems, " "))
	}

	if splitter != " " {
		printAble = printAble[1 : len(printAble)-1]
		if !strings.Contains(printAble, splitter) {
			splitter = " "
		}
	}

	for _, subStr := range strings.Split(printAble, splitter) {
		if subStr != "" {
			str := subStr

			// fmt.Printf("str: %s\n", str)

			if splitter == " " {
				printAbleTemp += strings.TrimSpace(str) + " "
			} else if strings.TrimSpace(str) == "" {
				printAbleTemp += " "
			} else {
				printAbleTemp += str
			}
		}
	}

	return strings.TrimSpace(printAbleTemp)
}

func Filter[T any](slice []T, condition func(T) bool) []T {
	var result []T
	for _, v := range slice {
		if condition(v) {
			result = append(result, v)
		}
	}
	return result
}

func parseDoubleQuotedStrings(input string) []string {
	var result []string
	var current strings.Builder
	inQuotes := false

	for i := 0; i < len(input); i++ {
		char := input[i]

		switch char {
		case '"':
			inQuotes = !inQuotes // Toggle quotes
		case ' ':
			if inQuotes {
				current.WriteByte(char) // Keep spaces inside quotes
			} else {
				if current.Len() > 0 {
					result = append(result, current.String())
					current.Reset()
				}
			}
		default:
			current.WriteByte(char)
		}
	}

	// Add any remaining characters to the result
	if current.Len() > 0 {
		result = append(result, current.String())
	}

	return result
}

func parseSingleQuotedStrings(input string) []string {
	var result []string
	var current string
	inQuotes := false

	for i := 0; i < len(input); i++ {
		char := input[i]

		switch char {
		case '\'':
			inQuotes = !inQuotes            // Toggle quote state
			if !inQuotes && current != "" { // If closing quote and current is not empty
				result = append(result, current)
				current = ""
			}
		case ' ':
			if inQuotes {
				current += string(char) // Preserve spaces within quotes
			} else if current != "" {
				result = append(result, current) // Add to result if not in quotes and not empty
				current = ""
			}
		default:
			current += string(char) // Build the current filename or text
		}
	}

	if current != "" { // Add any remaining text after loop ends
		result = append(result, current)
	}

	return result
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
