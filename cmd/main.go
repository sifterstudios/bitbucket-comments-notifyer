package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
)

func main() {
	http.HandleFunc("/webhook", func(w http.ResponseWriter, r *http.Request) {
		// Read the incoming webhoob payload
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			http.Error(w, "Error reading request body", http.StatusBadRequest)
			return
		}

		// Process the payload (e.g. validate, trigger actions, etc.)
		// You'll need to implement your Bitbucket event handling logic here

		fmt.Println("Received webhook payload:")
		fmt.Println(string(body))

		w.WriteHeader(http.StatusOK)
	})

	// Start the server on port 8080
	if err := http.ListenAndServe(":8080", nil); err != nil {
		fmt.Println("Error starting server:", err)
	}
}
