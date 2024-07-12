package api

import (
	"fmt"
	"os"
)

type Config struct {
	APIURL   string
	APIToken string
}

func LoadConfig() (*Config, error) {
	apiURL := os.Getenv("CAPACITIES_API_URL")
	if apiURL == "" {
		return nil, fmt.Errorf("CAPACITIES_API_URL environment variable is not set")
	}

	apiToken := os.Getenv("CAPACITIES_API_TOKEN")
	if apiToken == "" {
		return nil, fmt.Errorf("CAPACITIES_API_TOKEN environment variable is not set")
	}

	return &Config{
		APIURL:   apiURL,
		APIToken: apiToken,
	}, nil
}
