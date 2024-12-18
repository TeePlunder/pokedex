package api

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type Client struct {
	BaseUrl string
}

func NewClient(baseUrl string) *Client {
	return &Client{
		BaseUrl: baseUrl,
	}
}

func (c *Client) getResource(path string, v interface{}) error {
	url := fmt.Sprintf("%s/%s", c.BaseUrl, path)

	resp, err := http.Get(url)

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

func (c *Client) GetLocationAreas(path string) (LocationAreaResponse, error) {
	if path == "" {
		path = "location-area"
	}

	var res LocationAreaResponse

	if err := c.getResource(path, &res); err != nil {
		return res, err
	}

	return res, nil
}
