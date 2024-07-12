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

func FetchData(config *Config) {
	// Build the URL
	parsedURL, err := url.Parse(config.APIURL)
	if err != nil {
		fmt.Printf("Error parsing base URL: %v\n", err)
		return
	}
	parsedURL.Path = parsedURL.Path + "/spaces"
	url := parsedURL.String()

	// Create a new request
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		fmt.Printf("Error creating request: %v", err)
		return
	}

	// Add the token to the request headers
	req.Header.Add("Authorization", "Bearer "+config.APIToken)

	// Perform the request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Printf("Error sending request: %v", err)
		return
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("Error reading response body: %v", err)
		return
	}

	// Print the raw response body for debugging
	fmt.Println("Raw response body:", string(body))

	// Unmarshal the response body into a SpacesResponse struct
	var data SpacesResponse
	if err := json.Unmarshal(body, &data); err != nil {
		fmt.Printf("Error unmarshaling response body: %v", err)
		return
	}

	// Marshal the data back to a pretty-printed JSON string
	prettyJSON, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		fmt.Printf("Error marshaling JSON: %v", err)
		return
	}

	// Print the pretty-printed JSON
	fmt.Println(string(prettyJSON))
	// Update the model with the fetched data
}
