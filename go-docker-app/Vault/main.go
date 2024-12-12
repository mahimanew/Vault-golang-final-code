package Vault

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
)

// Define a struct to map the inner object
type APIData struct {
	ApiService string   `json:"ApiService"`
	Email      string   `json:"email"`
	Role       []string `json:"role"`
	UUID       string   `json:"uuid"`
}

func ApiKey() map[string]APIData {
	// Define the file path
	filePath := "/vault/secrets/app.txt"

	// Open the file
	file, err := os.Open(filePath)
	if err != nil {
		log.Fatalf("Failed to open file: %v", err)
	}
	defer file.Close()

	// Read the file content
	content, err := ioutil.ReadAll(file)
	if err != nil {
		log.Fatalf("Failed to read file: %v", err)
	}

	// Parse the JSON data into a map
	data := make(map[string]APIData)
	if err := json.Unmarshal(content, &data); err != nil {
		log.Fatalf("Failed to parse JSON: %v", err)
	}

	return data
}
