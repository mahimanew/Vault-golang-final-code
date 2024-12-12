package main

import (
	"fmt"
	vault "golang/go-docker-app/Vault"
	"net/http"
)

func handler(w http.ResponseWriter, r *http.Request) {

	fmt.Fprintf(w, "Hello, Dockerized Go!")
	data := vault.ApiKey()
	for _, value := range data {
		fmt.Fprintf(w, value.Email)
		fmt.Fprintf(w, value.ApiService)
	}

}

func main() {
	http.HandleFunc("/", handler)
	fmt.Println("Starting server on :8080...")
	http.ListenAndServe(":8080", nil)
}
