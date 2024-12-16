package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

type cliCommand struct {
	name        string
	description string
	callback    func() error
}

// The purpose of this function will be to split the users input into "words" based on whitespace.
// It should also lowercase the input and trim any leading or trailing whitespace
func cleanInput(text string) []string {
	text = strings.TrimSpace(text)
	return strings.Fields(text)
}

func commandExit() error {
	fmt.Println("Closing the Pokedex... Goodbye!")
	os.Exit(0)
	return nil
}

func main() {
	scanner := bufio.NewScanner(os.Stdin)

	commands := map[string]cliCommand{
		"exit": {
			name:        "exit",
			description: "Exit the Pokedex",
			callback:    commandExit,
		},
	}

	for {
		fmt.Print("Pokedex > ")

		if !scanner.Scan() {
			fmt.Fprintln(os.Stderr, "shouldn't see an error scanning a string")
		}

		if err := scanner.Err(); err != nil {
			fmt.Fprintln(os.Stderr, "shouldn't see an error scanning a string")
		}

		input := scanner.Text()
		words := cleanInput(input)
		firstWord := strings.ToLower(words[0])

		if cmd, exists := commands[firstWord]; exists {
			callbackErr := cmd.callback()
			if callbackErr != nil {
				fmt.Println("Error:", callbackErr)
			}
		} else {
			fmt.Println("Unknown command")
		}
	}
}
