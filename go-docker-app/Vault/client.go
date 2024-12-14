package Vault

import (
	"fmt"
	"log"
	"os"

	"github.com/hashicorp/vault/api"
)

// Define a struct to map the inner object
type APIData struct {
	ApiService string   `json:"ApiService"`
	Email      string   `json:"email"`
	Role       []string `json:"role"`
	UUID       string   `json:"uuid"`
}

func ApiKey() (string, error) {

	vaultAddress := os.Getenv("VAULT_ADDR")
	if vaultAddress == "" {
		log.Fatal("VAULT_ADDR environment variable is not set")
	}

	role := os.Getenv("VAULT_ROLE")
	if role == "" {
		log.Fatal("VAULT_ROLE environment variable is not set")
	}

	jwtPath := "/var/run/secrets/kubernetes.io/serviceaccount/token"

	// Read the JWT token
	jwtToken, err := os.ReadFile(jwtPath)
	if err != nil {
		log.Fatalf("Failed to read JWT token: %v", err)
	}

	// Create Vault client
	config := api.DefaultConfig()
	config.Address = vaultAddress

	client, err := api.NewClient(config)
	if err != nil {
		log.Fatalf("Failed to create Vault client: %v", err)
	}

	// Authenticate with Kubernetes Auth
	data := map[string]interface{}{
		"role": role,
		"jwt":  string(jwtToken),
	}

	authPath := "auth/kubernetes/login"
	secret, err := client.Logical().Write(authPath, data)
	if err != nil {
		log.Fatalf("Failed to authenticate with Vault: %v", err)
	}

	// Set the token for future requests
	client.SetToken(secret.Auth.ClientToken)

	// Read a secret from Vault
	secretPath := "kv-v2/fakeapp/mysecret"
	secret, err = client.Logical().Read(secretPath)
	if err != nil {
		log.Fatalf("Failed to read secret from Vault: %v", err)
	}

	if secret == nil || secret.Data["data"] == nil {
		log.Fatal("Secret not found or no data")
	}

	// Access the secret values
	secretData := secret.Data["data"].(map[string]interface{})

	result := ""
	for key, value := range secretData {
		result += fmt.Sprintf("%s: %s\n", key, value)
	}

	return result, nil
}
