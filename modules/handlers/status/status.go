package web_status

import (
	"GoRecordurbate/modules/bot"
	"GoRecordurbate/modules/config"
	"GoRecordurbate/modules/file"
	"encoding/json"
	"fmt"
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
	if r.Method != http.MethodGet {
		http.Error(w, "Only GET allowed", http.StatusMethodNotAllowed)
		return
	}

	var recorder Response
	config.Reload(file.Config_json_path, &config.Streamers)
	// Assuming config.Settings.App.Streamers is accessible
	for _, s := range bot.Bot.ListRecorders() {
		if s.IsRecording {
			recorder.Status = "Running"
			break // No need to continue checking
		} else {
			recorder.Status = "Stopped"

		}
	}
	recorder.BotStatus = bot.Bot.Status()
	fmt.Println(bot.Bot.Status())
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(recorder)
}
