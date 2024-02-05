package repl

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"

	"golang.org/x/exp/slices"
)

func YesOrNo(message string) bool {
	scanner := bufio.NewScanner(os.Stdin)
	fmt.Println(message)
	fmt.Println("(1) Yes  (2) No")
	for {
		fmt.Print("Enter number -> ")
		scanner.Scan()
		input := scanner.Text()
		switch input {
		case "1":
			return true
		case "2":
			return false
		default:
			printError("Please enter 1 or 2")
		}
	}
}

func Readline(message string, verifiers []func(string) error) string {
	scanner := bufio.NewScanner(os.Stdin)
	fmt.Println(message)
	for {
		fmt.Print("-> ")
		scanner.Scan()
		input := scanner.Text()
		input = strings.Trim(input, " ")
		if verifiers == nil {
			return input
		}
		ok := true
		for _, verify := range verifiers {
			err := verify(input)
			if err != nil {
				ok = false
				printError("Invalid input: " + err.Error())
			}
		}
		if ok {
			return input
		}
	}
}

func Selector(items []string, message string, multi bool) []string {
	scanner := bufio.NewScanner(os.Stdin)

	options := ""
	for i, item := range items {
		options += fmt.Sprintf("(%d) %s\t", i+1, item)
	}
	fmt.Println(options)
	fmt.Println(message)

	answer := []string{}
	for {
		fmt.Print("-> ")
		scanner.Scan()
		input := scanner.Text()

		noErr := true
		if len(input) == 0 {
			printError("Please select one")
			continue
		}
		selectedList := strings.Split(input, " ")
		if !multi && len(selectedList) > 1 {
			printError("Please select one")
			continue
		}
		for _, selected := range selectedList {
			num, err := strconv.Atoi(selected)
			if err != nil || num < 1 || num > len(items) {
				noErr = false
				printWarn("Invalid input: " + selected)
				continue
			}
			item := items[num-1]
			if !slices.Contains(answer, item) {
				answer = append(answer, item)
			}
		}
		if noErr {
			return answer
		} else {
			answer = []string{}
		}
	}
}
