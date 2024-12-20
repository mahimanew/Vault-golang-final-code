package main

import (
	"crypto/tls"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"github.com/hashicorp/vault/api"
)

func main() {
	// Read the Vault address and Vault role
	vaultAddr := os.Getenv("VAULT_ADDR")
	if vaultAddr == "" {
		log.Fatal("VAULT_ADDR environment variable is not set")
	}

	vaultRole := os.Getenv("VAULT_ROLE")
	if vaultRole == "" {
		log.Fatal("VAULT_ROLE environment variable is not set")
	}

	tokenPath := "/var/run/secrets/vault/token"

	// Read the token from the file
	token, err := ioutil.ReadFile(tokenPath)
	if err != nil {
		log.Fatalf("Error reading token: %v", err)
	}

	vaultToken := string(token)
	// Create a Vault client
	config := api.DefaultConfig()
	config.Address = vaultAddr

	client, err := api.NewClient(&api.Config{
		Address: vaultAddr,
		HttpClient: &http.Client{
			Transport: &http.Transport{
				TLSClientConfig: &tls.Config{
					InsecureSkipVerify: true, // Disable TLS verification for testing purposes
				},
			},
		},
	})
	if err != nil {
		log.Fatalf("Error creating Vault client: %v", err)
	}

	// Authenticate with Vault using Kubernetes Auth
	data := map[string]interface{}{
		"role": vaultRole,
		"jwt":  vaultToken,
	}

	authPath := "auth/kubernetes/login"
	secret, err := client.Logical().Write(authPath, data)
	if err != nil {
		log.Fatalf("Failed to authenticate with Vault: %v", err)
	}

	// Set the token for future requests
	client.SetToken(secret.Auth.ClientToken)

	// You can also print the Vault client token for debugging
	log.Printf("Vault Client Token: %s", secret.Auth.ClientToken)

	secretPath := "techassurance/data/data/dev" // Replace with your secret path
	secretData, err := client.Logical().Read(secretPath)
	if err != nil {
		log.Fatalf("failed to read secret from Vault: %v", err)
	}

	// Print the secret data in logs for verification
	log.Printf("Secret data: %v\n", secretData.Data["data"])

	// Serve the secret data in the browser
	http.HandleFunc("/secret", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")
		fmt.Fprintf(w, "<h1>Vault Secret Data:</h1>")
		fmt.Fprintf(w, "<pre>%v</pre>", secretData.Data["data"])
	})

	// Start the HTTP server
	fmt.Println("Starting server at :8080")
	err = http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatalf("Error starting server: %v", err)
	}
}
