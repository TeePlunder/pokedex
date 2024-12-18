package cli

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/teeplunder/pokedexcli/internal/api"
)

type config struct {
	Next     *string
	Previous *string
}

type cliCommand struct {
	name        string
	description string
	callback    func(*CLI) error
}

type CLI struct {
	commands map[string]cliCommand
	client   *api.Client
	config   *config
}

// The purpose of this function will be to split the users input into "words" based on whitespace.
// It should also lowercase the input and trim any leading or trailing whitespace
func cleanInput(text string) []string {
	text = strings.TrimSpace(text)
	return strings.Fields(text)
}

func commandExit(cli *CLI) error {
	fmt.Println("Closing the Pokedex... Goodbye!")
	os.Exit(0)
	return nil
}

func commandHelp(cli *CLI) error {
	fmt.Println("Welcome to the Pokedex!")
	fmt.Printf("Usage:\n\n")
	for _, cmd := range cli.commands {
		fmt.Printf("%s: %s\n", cmd.name, cmd.description)
	}
	return nil
}

func (cli *CLI) unknownCommand() error {
	fmt.Println("Unkown Command")
	return nil
}

func commandMap(cli *CLI) error {
	data, err := cli.client.GetLocationAreas()
	if err != nil {
		return fmt.Errorf("failed to get areas: %w", err)
	}

	for _, area := range data.Results {
		fmt.Println(area.Name)
	}

	cli.config.Next = &data.Next
	cli.config.Previous = data.Previous

	return nil
}

func NewCLI(client *api.Client) *CLI {
	return &CLI{
		client: client,
		commands: map[string]cliCommand{
			"map": {
				name:        "map",
				description: "Gets all location areas",
				callback:    commandMap,
			},
			"exit": {
				name:        "exit",
				description: "Exit the Pokedex",
				callback:    commandExit,
			},
			"help": {
				name:        "help",
				description: "Displays a help message",
				callback:    commandHelp,
			},
		},
		config: &config{
			Next:     nil,
			Previous: nil,
		},
	}
}

func (cli *CLI) Run() {
	scanner := bufio.NewScanner(os.Stdin)

	for {
		fmt.Print("Pokedex > ")

		if !scanner.Scan() {
			fmt.Fprintln(os.Stderr, "Error reading input")
			break
		}

		input := scanner.Text()
		words := cleanInput(input)

		if len(words) <= 0 {
			fmt.Println("Please enter a command. Type 'help' to see all available commands")
			continue
		}

		firstWord := strings.ToLower(words[0])

		if cmd, exists := cli.commands[firstWord]; exists {
			// Execute the command callback
			if err := cmd.callback(cli); err != nil {
				fmt.Println("Error:", err)
			}
		} else {
			cli.unknownCommand()
		}
	}
}
