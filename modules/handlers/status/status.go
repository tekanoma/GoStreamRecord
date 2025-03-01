package status

import (
	"GoRecordurbate/modules/bot"
	"GoRecordurbate/modules/db"
	"GoRecordurbate/modules/file"
	"GoRecordurbate/modules/handlers/cookies"
	"encoding/json"
	"net/http"
)

// Response is a generic response structure for our API endpoints.
type Response struct {
	Status    string        `json:"status"`
	Message   string        `json:"message"`
	Data      interface{}   `json:"data,omitempty"`
	BotStatus bot.BotStatus `json:"botStatus"`
}

func StatusHandler(w http.ResponseWriter, r *http.Request) {
	if !cookies.Session.IsLoggedIn(w, r) {
		http.Redirect(w, r, "/login", http.StatusFound)
		return
	}
	if r.Method != http.MethodGet {
		http.Error(w, "Only GET allowed", http.StatusMethodNotAllowed)
		return
	}

	var recorder Response
	db.Reload(file.API_keys_file, &db.Streamers)
	// Assuming db.Settings.App.Streamers is accessible
	for _, s := range bot.Bot.ListRecorders() {
		if s.IsRecording {
			recorder.Status = "Running"
			break // No need to continue checking
		} else {
			recorder.Status = "Stopped"

		}
	}
	recorder.BotStatus = bot.Bot.Status()

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(recorder)
}

func ResponseHandler(w http.ResponseWriter, r *http.Request, message string, data interface{}) {
	resp := Response{
		Status:  "success",
		Message: message,
		Data:    data,
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}
