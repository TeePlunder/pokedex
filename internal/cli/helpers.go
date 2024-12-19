package cli

import (
	"fmt"
	"strings"
)

func displayLocationAreas(cli *CLI, path string) error {
	data, err := cli.client.GetLocationAreas(path)
	if err != nil {
		return fmt.Errorf("failed to get areas: %w", err)
	}

	for _, area := range data.Results {
		fmt.Println(area.Name)
	}

	if data.Next == "" {
		cli.config.Next = nil
	} else {
		cli.config.Next = &data.Next
	}

	cli.config.Previous = data.Previous

	return nil
}

func getLastUrlPath(path string) string {
	lastSlashIndex := strings.LastIndex(path, "/")
	if lastSlashIndex == -1 {
		return ""
	}
	return path[lastSlashIndex+1:]
}

// The purpose of this function will be to split the users input into "words" based on whitespace.
// It should also lowercase the input and trim any leading or trailing whitespace
func CleanInput(text string) []string {
	text = strings.TrimSpace(text)
	text = strings.ToLower(text)
	return strings.Fields(text)
}
