package handlers

import (
	"GoRecordurbate/modules/bot"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
)

var Bot *bot.Bot

func Init() {

	Bot = bot.NewBot(log.New(os.Stdout, "lpg.log", log.LstdFlags))

}

// dcodes a JSON payload with a "command" field (start, stop, or restart)
// and returns a dummy response.
func controlHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Only POST allowed", http.StatusMethodNotAllowed)
		return
	}

	type RequestData struct {
		Command string `json:"command"`
	}
	fmt.Println("Got command")
	var reqData RequestData
	if err := json.NewDecoder(r.Body).Decode(&reqData); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}
	Bot.Command(reqData.Command)
	resp := Response{
		Message: fmt.Sprintf("Control command '%s' executed", reqData.Command),
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}
