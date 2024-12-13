package main

import (
	"fmt"
	"strings"
)

func main() {
	fmt.Println("Hello, World!")
}

// The purpose of this function will be to split the users input into "words" based on whitespace.
// It should also lowercase the input and trim any leading or trailing whitespace
func cleanInput(text string) []string {
	text = strings.TrimSpace(text)
	return strings.Fields(text)
}
