package web

import (
	"encoding/json"
	"fmt"
	"github.com/sifterstudios/bitbucket-comments-notifyer/internal/notification"
	"net/http"

	"github.com/gorilla/mux"
)

func StartWebServer() {
	r := mux.NewRouter()

	r.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("../web/static/"))))

	r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "../web/templates/index.html")
	})
	r.HandleFunc("/send-notification", sendNotificationHandler).Methods("POST")
	r.HandleFunc("/manual-update", manualUpdateHandler).Methods("POST")

	fmt.Println("Listening on port 8080")
	http.ListenAndServe(":8080", r)
}

func manualUpdateHandler(writer http.ResponseWriter, request *http.Request) {

}

func sendNotificationHandler(writer http.ResponseWriter, request *http.Request) {
	err := notification.SendNotification("Test notification", "It just works! :D")

	if err != nil {
		http.Error(writer, "Failed to send notification", http.StatusInternalServerError)
		return
	}

	response := map[string]string{
		"message": "Notification sent successfully, look top right!",
	}
	writer.Header().Set("Content-Type", "application/json")
	json.NewEncoder(writer).Encode(response)

}
