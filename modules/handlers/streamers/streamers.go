package streamers

import (
	"GoRecordurbate/modules/bot"
	"GoRecordurbate/modules/config"
	"GoRecordurbate/modules/handlers/status"
	"encoding/json"
	"net/http"
)

// Handles POST /api/add-streamer.
// It decodes a JSON payload with a "data" field and returns a dummy response.
func AddStreamer(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Only POST allowed", http.StatusMethodNotAllowed)
		return
	}

	type RequestData struct {
		Data string `json:"data"`
	}
	var reqData RequestData
	if err := json.NewDecoder(r.Body).Decode(&reqData); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	resp := status.Response{
		Message: config.AddStreamer(reqData.Data),
		Data:    reqData.Data,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

// Handles POST /api/remove-streamer.
// It decodes a JSON payload with the selected option and returns a dummy response.
func RemoveStreamer(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Only POST allowed", http.StatusMethodNotAllowed)
		return
	}

	type RequestData struct {
		Selected string `json:"selected"`
	}
	var reqData RequestData
	if err := json.NewDecoder(r.Body).Decode(&reqData); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	resp := status.Response{
		Message: config.RemoveStreamer(reqData.Selected),
		Data:    reqData.Selected,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

// Handles GET /api/get-streamers.
func GetStreamers(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Only GET allowed", http.StatusMethodNotAllowed)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	list := []string{}
	for _, s := range config.Streamers.StreamerList {
		list = append(list, s.Name)
	}
	json.NewEncoder(w).Encode(list)
}

// Handles POST /api/stop-singleProcess.
// It decodes a JSON payload with the selected option and returns a dummy response.
func StopProcess(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Only POST allowed", http.StatusMethodNotAllowed)
		return
	}

	type RequestData struct {
		Streamer string `json:"streamer"`
	}
	var reqData RequestData
	if err := json.NewDecoder(r.Body).Decode(&reqData); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}
	status.ResponseHandler(w, r, "Stopping process for "+reqData.Streamer, nil)
	bot.Bot.StopProcess(reqData.Streamer)
	status.ResponseHandler(w, r, "Stopped process for"+reqData.Streamer, nil)
}
