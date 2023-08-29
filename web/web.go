package web

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"reflect"
	"strconv"
	"sync"
	"time"

	"github.com/sifterstudios/bitbucket-notifier/bitbucket"
	"github.com/sifterstudios/bitbucket-notifier/data"
	"github.com/sifterstudios/bitbucket-notifier/notification"
)

func StartWebServer() {
	r := mux.NewRouter()

	r.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("../web/static/"))))
	r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "../web/templates/index.html")
	})
	r.HandleFunc("/send-notification", sendNotificationHandler).Methods("POST")
	r.HandleFunc("/update", updateHandler).Methods("GET")
	r.HandleFunc("/stats", getStatsHandler).Methods("GET")
	r.HandleFunc("/config", getConfigHandler).Methods("GET")
	r.HandleFunc("/config", setConfigHandler).Methods("POST")

	fmt.Println("Listening on port 1337")
	fmt.Println("Go to http://localhost:1337 to change settings and test the setup!")
	go func() {
		err := http.ListenAndServe(":1337", r)
		if err != nil {
			fmt.Println(err)
		}
	}()
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		startScheduledUpdate()
		wg.Done()
	}()
	wg.Wait()
}

func startScheduledUpdate() {
	for {
		println("Updating..." + time.Now().String())
		updateHandler(nil, nil)
		time.Sleep(time.Duration(data.UserConfig.ConfigNotifications.PollingInterval) * time.Minute)
	}
}
func getStatsHandler(writer http.ResponseWriter, _ *http.Request) {
	if data.CurrentPrs == nil {
		fmt.Println("Error: Could not find any Pull Requests")
		return
	}
	stats := data.ConvertActivePrResponseToUiStatistics(data.CurrentPrs)

	writer.Header().Set("Content-Type", "application/json")
	_, err := writer.Write([]byte(fmt.Sprintf(
		`{"lastUpdate": %d, "numberOfActivePrComments": %d, "numberOfActivePrTasks": %d}`,
		stats.LastUpdate, stats.NumberOfActivePrComments, stats.NumberOfActivePrTasks)))
	if err != nil {
		return
	}
}

func getConfigHandler(writer http.ResponseWriter, _ *http.Request) {
	config := data.UserConfig.ConfigNotifications
	fmt.Println(config)
	jsonUiStats, err := json.Marshal(config)
	if err != nil {
		fmt.Println(err)
	}
	_, err = writer.Write(jsonUiStats)
	if err != nil {
		return
	}
}
func setConfigHandler(_ http.ResponseWriter, request *http.Request) {
	err := request.ParseForm()
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(request.Form)
	newConfig := data.ConfigNotifications{}
	pollingInterval, err := strconv.Atoi(request.Form.Get("monitoringFrequencyInput"))
	if err != nil {
		fmt.Println(err)
	}
	newConfig.PollingInterval = pollingInterval
	newConfig.Comments = request.Form.Get("notifyCommentsCheckbox") == "on"
	newConfig.Tasks = request.Form.Get("notifyTasksCheckbox") == "on"
	newConfig.StatusChanges = request.Form.Get("notifyStatusChangesCheckbox") == "on"
	newConfig.CompletionTime = request.Form.Get("notifyCompletionTimeCheckbox") == "on"

	data.UserConfig.ConfigNotifications = newConfig
}

func updateHandler(writer http.ResponseWriter, _ *http.Request) {
	prResponse, err := bitbucket.GetCurrentPullRequestsByUser(data.UserConfig)
	if err != nil {
		log.Print(err)
	}

	if !reflect.DeepEqual(prResponse.Values, data.CurrentPrs) {
		data.HandleCurrentPrs(prResponse.Values)
	}

	activityResponse, err := bitbucket.GetPullRequestsActivity(data.CurrentPrs)
	data.HandlePrActivity(data.CurrentPrs, activityResponse)

	uiStats := data.ConvertActivePrResponseToUiStatistics(prResponse.Values)

	fmt.Println(uiStats)

	if writer != nil {
		jsonUiStats, err := json.Marshal(uiStats)
		_, err = writer.Write(jsonUiStats)
		if err != nil {
			return
		}
	}
}

func sendNotificationHandler(writer http.ResponseWriter, _ *http.Request) {
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
