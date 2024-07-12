package api

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
)

func FetchSpacesData(config *Config) ([]byte, error) {
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
