package web

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
)

func StartWebServer() {
	r := mux.NewRouter()

	r.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("../web/static/"))))

	r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "../web/templates/index.html")
	})

	fmt.Println("Listening on port 8080")
	http.ListenAndServe(":8080", r)
}
