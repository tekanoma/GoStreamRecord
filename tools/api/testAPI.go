package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
)

func main() {
	// Define command-line flags.
	apiKey := flag.String("api_key", "", "API key to use for the request")
	ip := flag.String("ip", "localhost", "IP address of the server")
	port := flag.Int("port", 8055, "Port of the server")
	flag.Parse()

	// Construct the URL for the health endpoint.
	url := fmt.Sprintf("http://%s:%d/api/v1/health", *ip, *port)
	log.Printf("Sending request to %s", url)

	// Create a new GET request.
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Fatalf("Error creating request: %v", err)
	}

	// If an API key was provided, add it to the request header.
	if *apiKey != "" {
		req.Header.Add("X-API-Key", *apiKey)
	}

	// Send the request.
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Fatalf("Error sending request: %v", err)
	}
	defer resp.Body.Close()

	// Read the response body.
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatalf("Error reading response: %v", err)
	}

	// Output the response.
	fmt.Printf("Response Status: %s\n", resp.Status)
	fmt.Printf("Response Body: %s\n", body)
}
