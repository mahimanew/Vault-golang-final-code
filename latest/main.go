package main

import (
	"crypto/tls"
	"fmt"
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

	vaultToken := os.Getenv("VAULT_SA_TOKEN")
	//vaultToken := "eyJhbGciOiJSUzI1NiIsImtpZCI6InlYb1QxbmhZS2pnYTdOR3dxenBKMU8wRWJlMmhvRUVmYXpWZWxQbEJnbVEifQ.eyJpc3MiOiJrdWJlcm5ldGVzL3NlcnZpY2VhY2NvdW50Iiwia3ViZXJuZXRlcy5pby9zZXJ2aWNlYWNjb3VudC9uYW1lc3BhY2UiOiJ2YXVsdCIsImt1YmVybmV0ZXMuaW8vc2VydmljZWFjY291bnQvc2VjcmV0Lm5hbWUiOiJ2YXVsdC1zYS10b2tlbiIsImt1YmVybmV0ZXMuaW8vc2VydmljZWFjY291bnQvc2VydmljZS1hY2NvdW50Lm5hbWUiOiJ2YXVsdC1zYSIsImt1YmVybmV0ZXMuaW8vc2VydmljZWFjY291bnQvc2VydmljZS1hY2NvdW50LnVpZCI6IjBmNjllYzgzLTljODYtNDRhYy1iYzBkLWNiNDhiN2FmNTEyMiIsInN1YiI6InN5c3RlbTpzZXJ2aWNlYWNjb3VudDp2YXVsdDp2YXVsdC1zYSJ9.ebaAYjDMz7lWq_LI1B6yzx8hA-1ac8gRDkSwfmKjschfiQsW72fyDD6wBVPDSWJrHx0H85ethOTZwVfj86Rwh2rPy5ZlCMtw-ezGwuHW5KJ9Tu7TE3lU9ZESsWWhVko-jDPMNHMBuaUxe3jfOL9TsrlY_PBES_lXuF0Ds_WeAW_MK3Md0Lw3MWTBWIPY1i7e97C3nC_7cxHfItJlyuNcLIAXDpIjb3eaAJAJig5PURJwP7d5w-J6JtpaG8SFbnACzsz-VSeZ5Hxn2IfdacHIPExIVST_RQTVA3o3KYOPORJmR7XeLgD9uhQGZABabTSi2acAzV-IEIRcloy3RltHFA"
	if vaultToken == "" {
		log.Fatal("VAULT_SA_TOKEN environment variable is not set")
	}

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
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
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
