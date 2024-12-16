package main

import (
	"github.com/teeplunder/pokedexcli/internal/api"
	"github.com/teeplunder/pokedexcli/internal/cli"
)

func main() {
	client := api.NewClient("https://pokeapi.co/api/v2/")
	app := cli.NewCLI(client)
	app.Run()
}
