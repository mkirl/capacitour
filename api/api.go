package api

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
)

type Icon struct {
	Type string `json:"type"`
	Val  string `json:"val"`
}

type Space struct {
	ID    string `json:"id"`
	Title string `json:"title"`
	Icon  Icon   `json:"icon"`
}

type SpacesResponse struct {
	Spaces []Space `json:"spaces"`
}

func FetchAllSpacesData(config *Config) ([]byte, error) {
	// Build the URL
	parsedURL, err := url.Parse(config.APIURL)
	if err != nil {
		fmt.Printf("Error parsing base URL: %v\n", err)
		return nil, err
	}
	parsedURL.Path = parsedURL.Path + "/spaces"
	url := parsedURL.String()

	// Create a new request
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		fmt.Printf("Error creating request: %v", err)
		return nil, err
	}

	// Add the token to the request headers
	req.Header.Add("Authorization", "Bearer "+config.APIToken)

	// Perform the request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Printf("Error sending request: %v", err)
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("Error reading response body: %v", err)
		return nil, err
	}

	// Print the raw response body for debugging
	// Marshal the data back to a pretty-printed JSON string
	return body, nil
}

func FetchAllSpaces(config *Config) ([]Space, error) {
	body, err := FetchAllSpacesData(config)
	if err != nil {
		return nil, err
	}

	var spacesResponse SpacesResponse
	err = json.Unmarshal(body, &spacesResponse)
	if err != nil {
		return nil, err
	}

	return spacesResponse.Spaces, nil
}

func FetchSpaceData(config *Config, space Space) ([]byte, error) {
	parsedURL, err := url.Parse(config.APIURL)
	if err != nil {
		return nil, err
	}
	parsedURL.Path = parsedURL.Path + "/spaces-info"
	query := parsedURL.Query()
	query.Set("spaceid", space.ID)
	parsedURL.RawQuery = query.Encode()
	url := parsedURL.String()

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Add("Authorization", "Bearer "+config.APIToken)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	return body, nil
}

func FetchSpace(config *Config, space Space) ([]byte, error) {
	body, err := FetchSpaceData(config, space)
	if err != nil {
		return nil, err
	}
	return body, nil
}
