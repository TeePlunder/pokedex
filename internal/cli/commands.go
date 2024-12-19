package cli

import (
	"bufio"
	"fmt"
	"os"

	"github.com/teeplunder/pokedexcli/internal/api"
)

type config struct {
	Next     *string
	Previous *string
}

type cliCommand struct {
	name        string
	description string
	callback    func(*CLI, string) error
}

type CLI struct {
	commands map[string]cliCommand
	client   *api.Client
	config   *config
}

func commandExit(cli *CLI, param string) error {
	fmt.Println("Closing the Pokedex... Goodbye!")
	os.Exit(0)
	return nil
}

func commandHelp(cli *CLI, param string) error {
	fmt.Println("Welcome to the Pokedex!")
	fmt.Println("Usage:")
	fmt.Println("")
	for _, cmd := range cli.commands {
		fmt.Printf("%s: %s\n", cmd.name, cmd.description)
	}
	return nil
}

func (cli *CLI) unknownCommand() error {
	fmt.Println("Unknown Command")
	return nil
}

func commandMap(cli *CLI, param string) error {
	path := ""
	if cli.config.Next != nil {
		// get url path (last part from the url)
		path = getLastUrlPath(*cli.config.Next)
	}

	return displayLocationAreas(cli, path)
}

func commandMapBack(cli *CLI, param string) error {
	if cli.config.Previous == nil {
		fmt.Println("you're on the first page")
		return nil
	}
	path := getLastUrlPath(*cli.config.Previous)

	return displayLocationAreas(cli, path)
}

func commandExplore(cli *CLI, param string) error {
	area := param
	fmt.Printf("Exploring %s...\n", area)

	data, err := cli.client.GetPokemonEncountersAtLocationArea(area)
	if err != nil {
		return fmt.Errorf("failed to explore %s: %w", area, err)
	}

	fmt.Println("Found Pokemon:")
	for _, pokemon := range data {
		fmt.Printf("- %s\n", pokemon.Name)
	}

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
			"mapb": {
				name:        "mapb",
				description: "Gets all location areas from the previous location",
				callback:    commandMapBack,
			},
			"explore": {
				name:        "explore <area>",
				description: "Explores an area and displays all Pokemon found",
				callback:    commandExplore,
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
		words := CleanInput(input)

		if len(words) == 0 {
			fmt.Println("Please enter a command. Type 'help' to see all available commands")
			continue
		}

		if cmd, exists := cli.commands[words[0]]; exists {
			param := ""
			if len(words) > 1 {
				param = words[1]
			}

			// Execute the command callback
			if err := cmd.callback(cli, param); err != nil {
				fmt.Println("Error:", err)
			}
		} else {
			cli.unknownCommand()
		}
	}
}
