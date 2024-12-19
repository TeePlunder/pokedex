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

// fetches the resource at the given path and unmarshals the JSON response into v.
func (c *Client) getResource(urlPath string, v interface{}) error {
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

	if err := json.NewDecoder(resp.Body).Decode(v); err != nil {
		return fmt.Errorf("Failed to decode json: %w", err)
	}

	return nil
}

// Retrieves location areas from the given path.
// If path is empty, it defaults to "location-area".
func (c *Client) GetLocationAreas(path string) (LocationAreaResponse, error) {
	if path == "" {
		path = "location-area"
	}

	var res LocationAreaResponse

	// check if value exists in cache

	if cachedValue, ok := c.Cache.Get(path); ok {
		// cachedValue is []byte (raw JSON)
		// decode into res
		if err := json.Unmarshal(cachedValue, &res); err == nil {
			fmt.Println("-> use cached values")
			return res, nil
		}
		// fall throught to api call
	}

	if err := c.getResource(path, &res); err != nil {
		return res, err
	}

	// store value in Cache

	// convert res to json bytes for cache storage
	jsonData, err := json.Marshal(res)
	if err == nil {
		c.Cache.Add(path, jsonData)
	}

	fmt.Println("-> use api values")
	return res, nil
}
