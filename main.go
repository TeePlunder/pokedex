package main

import (
	"time"

	"github.com/teeplunder/pokedexcli/internal/api"
	"github.com/teeplunder/pokedexcli/internal/cache"
	"github.com/teeplunder/pokedexcli/internal/cli"
)

func main() {
	cache := cache.NewCache(5 * time.Second)
	client := api.NewClient(api.API_BASE_PATH, cache)
	app := cli.NewCLI(client)
	app.Run()
}
