package cookies

import (
	"GoRecordurbate/modules/db"
	"GoRecordurbate/modules/file"
	"encoding/json"
	"fmt"
	"net/http"
)

// API key generation response
type api_response struct {
	Status  bool   `json:"status"`
	Message string `json:"message"`
	Key     string `json:"key"`
}

func GenAPIKeyHandler(w http.ResponseWriter, r *http.Request) {
	if !Session.IsLoggedIn(w, r) {
		http.Redirect(w, r, "/login", http.StatusFound)
		return
	}

	err := db.API.Load()
	if err != nil {
		fmt.Println(err)
		http.Error(w, "Error getting existing keys..", http.StatusBadRequest)
		return
	}

	session, err := Session.Store().Get(r, "session")
	new_api_config := db.API.NewKey()
	new_api_config.User = session.Values["user"].(string)

	if new_api_config.User == "" {
		http.Error(w, "Unable generate api keys..", http.StatusForbidden)
		return
	}

	new_api_config.Name = r.URL.Query().Get("name")

	for _, k := range db.API.Keys {
		if k.Name == new_api_config.Name {
			if err != nil {
				fmt.Println(err)
				http.Error(w, "Named key already exists!", http.StatusConflict)
				return
			}
		}
	}

	key, err := new_api_config.GenerateAPIKey(32)
	if err != nil {
		fmt.Println(err)
		http.Error(w, "Unable generate api keys..", http.StatusBadRequest)
		return
	}
	hashedKey, err := new_api_config.HashAPIKey(key)

	new_api_config.Key = hashedKey

	db.API.Keys = append(db.API.Keys, new_api_config)
	err = file.WriteJson(file.API_keys_file, db.API)
	if err != nil {
		fmt.Println(err)
		http.Error(w, "error saving new key..", http.StatusBadRequest)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	response := api_response{Status: true, Message: "Generated api key.", Key: key}
	json.NewEncoder(w).Encode(response)
}

func GetAPIkeys(w http.ResponseWriter, r *http.Request) {
	if !Session.IsLoggedIn(w, r) {
		http.Redirect(w, r, "/login", http.StatusFound)
		return
	}

	err := db.API.Load()
	if err != nil {
		fmt.Println(err)
		http.Error(w, "Error retrieving API keys: "+err.Error(), http.StatusInternalServerError)
		return
	}

	type data struct {
		Name string `json:"name"` // The field should start with an uppercase letter
	}
	var apiList []data
	for _, k := range db.API.Keys {
		apiList = append(apiList, data{Name: k.Name})
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK) // Explicitly set status code for successful response
	json.NewEncoder(w).Encode(apiList)
}

func DeleteAPIKeyHandler(w http.ResponseWriter, r *http.Request) {
	if !Session.IsLoggedIn(w, r) {
		http.Redirect(w, r, "/login", http.StatusFound)
		return
	}

	var tmp_secrets db.API_secrets

	type data struct {
		Name string `json:"new"`
	}
	type request struct {
		Data data `json:"data"`
	}
	var reqData request

	if err := json.NewDecoder(r.Body).Decode(&reqData); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	err := db.API.Load()
	if err != nil {
		fmt.Println(err)
		http.Error(w, "Error getting existing keys..", http.StatusBadRequest)
		return
	}

	session, err := Session.Store().Get(r, "session")
	username := session.Values["user"].(string)

	for _, k := range db.API.Keys {
		if k.Name == reqData.Data.Name && k.User == username {
			continue
		}
		tmp_secrets.Keys = append(tmp_secrets.Keys, k)
	}

	db.API.Keys = tmp_secrets.Keys

	err = file.WriteJson(file.API_keys_file, db.API)
	if err != nil {
		fmt.Println(err)
		http.Error(w, "error saving new key..", http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	response := api_response{Status: true, Message: "Deleted api key.", Key: "nil"}
	json.NewEncoder(w).Encode(response)
}
