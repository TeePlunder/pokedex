package cli

import (
	"bufio"
	"fmt"
	"math/rand/v2"
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
	pokedex  map[string]api.Pokemon
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
	if len(param) == 0 {
		return fmt.Errorf("Please enter the name of the area")
	}
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

func commandCatch(cli *CLI, param string) error {
	if len(param) == 0 {
		return fmt.Errorf("Please enter the name of the pokemon")
	}
	pokemon := param
	fmt.Printf("Throwing a Pokeball at %s...\n", pokemon)

	data, err := cli.client.GetPokemon(pokemon)
	if err != nil {
		return fmt.Errorf("failed to get pokemon %s: %w\n", pokemon, err)
	}

	max := data.BaseExperience
	catchValue := rand.IntN(max)
	catched := catchValue >= max/2

	if catched {
		fmt.Printf("%s was caught!\n", pokemon)
		cli.pokedex[pokemon] = api.Pokemon{
			Name: pokemon,
			URL:  "",
		}
	} else {
		fmt.Printf("%s escaped!\n", pokemon)
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
			"catch": {
				name:        "catch <pokemon>",
				description: "Tries to catch a pokemon",
				callback:    commandCatch,
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
		pokedex: make(map[string]api.Pokemon),
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
