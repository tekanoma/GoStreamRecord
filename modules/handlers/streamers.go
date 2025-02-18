package handlers

import (
	"GoRecordurbate/modules/config"
	"encoding/json"
	"net/http"
)

// sendStringHandler handles POST /api/add-streamer.
// It decodes a JSON payload with a "data" field and returns a dummy response.
func addStreamer(w http.ResponseWriter, r *http.Request) {
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

	resp := Response{
		Message: config.C.App.AddStreamer(reqData.Data),
		Data:    reqData.Data,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

// getListHandler handles GET /api/get-streamers.
// It returns a dummy list of options.
func getStreamers(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Only GET allowed", http.StatusMethodNotAllowed)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	list := []string{}
	for _, s := range config.C.App.Streamers {
		list = append(list, s.Name)
	}
	json.NewEncoder(w).Encode(list)
}

// selectItemHandler handles POST /api/remove-streamer.
// It decodes a JSON payload with the selected option and returns a dummy response.
func removeStreamer(w http.ResponseWriter, r *http.Request) {
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

	resp := Response{
		Message: config.C.App.RemoveStreamer(reqData.Selected),
		Data:    reqData.Selected,
	}
	
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}
