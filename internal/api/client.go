package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"

	"github.com/teeplunder/pokedexcli/internal/cache"
)

type Client struct {
	BaseUrl string
	Cache   *cache.Cache
}

func NewClient(baseUrl string, cache *cache.Cache) *Client {
	return &Client{
		BaseUrl: baseUrl,
		Cache:   cache,
	}
}

// Fetches the resource at the given path, attempts to load from cache first,
// and if not found, fetches from the API and caches the result.
func (c *Client) getResource(urlPath string, value interface{}) error {

	if cachedValue, ok := c.Cache.Get(urlPath); ok {
		// cachedValue is []byte (raw JSON)
		// decode into value
		if err := json.Unmarshal(cachedValue, value); err == nil {
			fmt.Println("-> use cached values")
			return nil
		}
		// fall throught to api call
	}

	base, err := url.Parse(c.BaseUrl)
	if err != nil {
		return fmt.Errorf("failed to parse base URL: %w", err)
	}

	paramUrl, err := url.Parse(urlPath)
	if err != nil {
		return fmt.Errorf("failed to parse relative URL: %w", err)
	}

	// merges base url and param url correctly together (with query params)
	reqUrl := base.ResolveReference(paramUrl)

	resp, err := http.Get(reqUrl.String())

	if err != nil {
		return fmt.Errorf("failed to perform GET request\n %w", err)
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	if err := json.NewDecoder(resp.Body).Decode(value); err != nil {
		return fmt.Errorf("Failed to decode json: %w", err)
	}

	// store value in Cache

	// convert value to json bytes for cache storage
	if jsonData, err := json.Marshal(value); err == nil {
		c.Cache.Add(urlPath, jsonData)
	}

	fmt.Println("-> use api values")
	return nil
}

// Retrieves location areas from the given path.
// If path is empty, it defaults to "location-area".
func (c *Client) GetLocationAreas(path string) (LocationAreaResponse, error) {
	if path == "" {
		path = API_PATH_LOCATION_AREA
	}

	var res LocationAreaResponse

	if err := c.getResource(path, &res); err != nil {
		return res, err
	}

	return res, nil
}

func (c *Client) GetPokemonEncountersAtLocationArea(area string) ([]Pokemon, error) {
	var res LocationAreaDetailsResponse
	var encounteredPokemons []Pokemon
	path := fmt.Sprintf(API_PATH_LOCATION_AREA_DETAILS, area)

	if err := c.getResource(path, &res); err != nil {
		return nil, err
	}

	for _, encounter := range res.PokemonEncounters {
		encounteredPokemons = append(encounteredPokemons, encounter.Pokemon)
	}

	return encounteredPokemons, nil

}
