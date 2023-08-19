package main

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
)

func main() {
	r := mux.NewRouter()

	r.PathPrefix("/static").Handler(http.StripPrefix("/static", http.FileServer(http.Dir("static"))))

	r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "templates/config.html")
	})

	port := 8080
	fmt.Printf("Listening on port %d\n", port)
	http.ListenAndServe(fmt.Sprintf(":%d", port), r)
}
