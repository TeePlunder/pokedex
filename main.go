package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

// The purpose of this function will be to split the users input into "words" based on whitespace.
// It should also lowercase the input and trim any leading or trailing whitespace
func cleanInput(text string) []string {
	text = strings.TrimSpace(text)
	return strings.Fields(text)
}

func main() {
	scanner := bufio.NewScanner(os.Stdin)

	for {
		fmt.Print("Pokedex > ")

		if !scanner.Scan() {
			fmt.Fprintln(os.Stderr, "shouldn't see an error scanning a string")
		}

		if err := scanner.Err(); err != nil {
			fmt.Fprintln(os.Stderr, "shouldn't see an error scanning a string")
		}

		cmd := scanner.Text()
		words := cleanInput(cmd)
		first := strings.ToLower(words[0])
		fmt.Printf("Your command was: %s\n", first)
	}
}
