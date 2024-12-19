package main

import (
	"time"

	"github.com/teeplunder/pokedexcli/internal/api"
	"github.com/teeplunder/pokedexcli/internal/cache"
	"github.com/teeplunder/pokedexcli/internal/cli"
)

func main() {
	cache := cache.NewCache(5 * time.Second)
	client := api.NewClient("https://pokeapi.co/api/v2/", cache)
	app := cli.NewCLI(client)
	app.Run()
}
