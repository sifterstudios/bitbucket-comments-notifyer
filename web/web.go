package web

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"

	"github.com/sifterstudios/bitbucket-comments-notifyer/data"
	"github.com/sifterstudios/bitbucket-comments-notifyer/internal/bitbucket"
	"github.com/sifterstudios/bitbucket-comments-notifyer/internal/notification"
)

func StartWebServer() {
	r := mux.NewRouter()

	r.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("../web/static/"))))

	r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "../web/templates/index.html")
	})
	r.HandleFunc("/send-notification", sendNotificationHandler).Methods("POST")
	r.HandleFunc("/manual-update", manualUpdateHandler).Methods("GET")
	r.HandleFunc("/stats", getStatsHandler).Methods("GET")
	r.HandleFunc("/config", getConfig).Methods("GET")

	fmt.Println("Listening on port 8080")
	http.ListenAndServe(":8080", r)
}

func getStatsHandler(writer http.ResponseWriter, request *http.Request) {
	response, err := bitbucket.GetActivePullRequestsByUser(data.UserConfig)
	if err != nil {
		log.Print(err)
	}
	uiStats := data.ConvertActivePrResponseToUiStatistics(response)
	fmt.Println(uiStats)
	jsonUiStats, err := json.Marshal(uiStats)
	writer.Write(jsonUiStats)
}

func getConfig(writer http.ResponseWriter, request *http.Request) {
	config := data.UserConfig.Notification
	fmt.Println(config)
	jsonUiStats, err := json.Marshal(config)
	if err != nil {
		fmt.Println(err)
	}
	writer.Write(jsonUiStats)
}

func manualUpdateHandler(writer http.ResponseWriter, request *http.Request) {
	response, err := bitbucket.GetActivePullRequestsByUser(data.UserConfig)
	if err != nil {
		log.Print(err)
	}
	uiStats := data.ConvertActivePrResponseToUiStatistics(response)
	fmt.Println(uiStats)
	jsonUiStats, err := json.Marshal(uiStats)
	writer.Write(jsonUiStats)
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
	err = json.NewEncoder(writer).Encode(response)
	if err != nil {
		return
	}
}
