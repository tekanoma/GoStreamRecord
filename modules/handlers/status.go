package handlers

import (
	"GoRecordurbate/modules/config"
	"encoding/json"
	"net/http"
)

func statusHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Only GET allowed", http.StatusMethodNotAllowed)
		return
	}

	type Rec struct {
		Status string `json:"status"`
	}

	var recorder Rec
	config.C.Reload()
	// Assuming config.C.App.Streamers is accessible
	for _, s := range Bot.ListRecorders() {
		if s.Running {
			recorder.Status = "Running"
			break // No need to continue checking
		} else {
			recorder.Status = "Stopped"

		}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(recorder)
}
