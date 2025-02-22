package web_config

import (
	"GoRecordurbate/modules/config"
	"GoRecordurbate/modules/file"
	"GoRecordurbate/modules/handlers/cookies"
	"GoRecordurbate/modules/handlers/login"
	web_response "GoRecordurbate/modules/handlers/response"
	"encoding/json"
	"net/http"
)

// sendStringHandler handles POST /api/add-streamer.
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

	resp := web_response.Response{
		Message: config.AddStreamer(reqData.Data),
		Data:    reqData.Data,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

// getListHandler handles GET /api/get-streamers.
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
func GetUsers(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Only GET allowed", http.StatusMethodNotAllowed)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	json.NewEncoder(w).Encode(config.Users.Users)
}
func UpdateUsers(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Only POST allowed", http.StatusMethodNotAllowed)
		return
	}

	type RequestData struct {
		OldUsername string `json:"oldUsername"`
		NewUsername string `json:"newUsername"`
		NewPassword string `json:"newPassword"`
	}
	var reqData RequestData
	if err := json.NewDecoder(r.Body).Decode(&reqData); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}
	modified := false
	for i, user := range config.Users.Users {
		if user.Name == reqData.OldUsername {
			config.Users.Users[i].Name = reqData.NewUsername
			config.Users.Users[i].Key = string(login.HashedPassword(reqData.NewPassword))
			modified = true
			break
		}
	}
	if modified {
		config.Update(file.Users_json_path, config.Users)
	}

	resp := web_response.Response{
		Message: "User modified!",
	}
	for _, u := range config.Users.Users {
		cookies.UserStore[u.Name] = u.Key
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)

}

func AddUser(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Only POST allowed", http.StatusMethodNotAllowed)
		return
	}

	type RequestData struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}
	var reqData RequestData
	if err := json.NewDecoder(r.Body).Decode(&reqData); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}
	for _, user := range config.Users.Users {
		if user.Name == reqData.Username {
			resp := web_response.Response{
				Message: "User already exists!",
			}
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(resp)
			return
		}
	}

	config.Users.Users = append(config.Users.Users, config.Login{Name: reqData.Username, Key: string(login.HashedPassword(reqData.Password))})
	config.Update(file.Users_json_path, config.Users)

	resp := web_response.Response{
		Message: reqData.Username + " added!",
	}
	for _, u := range config.Users.Users {
		cookies.UserStore[u.Name] = u.Key
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)

}

// selectItemHandler handles POST /api/remove-streamer.
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

	resp := web_response.Response{
		Message: config.RemoveStreamer(reqData.Selected),
		Data:    reqData.Selected,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}
