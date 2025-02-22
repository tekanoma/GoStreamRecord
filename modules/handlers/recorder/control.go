package web_recorder

import (
	"GoRecordurbate/modules/bot"
	web_status "GoRecordurbate/modules/handlers/status"
	"encoding/json"
	"fmt"
	"net/http"
)

// dcodes a JSON payload with a "command" field (start, stop, or restart)
// and returns a dummy response.
func ControlHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Only POST allowed", http.StatusMethodNotAllowed)
		return
	}

	type RequestData struct {
		Command string `json:"command"`
	}
	var reqData RequestData
	if err := json.NewDecoder(r.Body).Decode(&reqData); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	bot.Bot.Command(reqData.Command)
	resp := web_status.Response{
		Message: fmt.Sprintf("Control command '%s' executed", reqData.Command),
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)

}
