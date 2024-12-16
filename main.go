package main

import (
	"github.com/teeplunder/pokedexcli/internal/cli"
)

func main() {
	app := cli.NewCLI()
	app.Run()
}
