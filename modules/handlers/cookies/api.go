package cookies

import (
	"GoRecordurbate/modules/file"
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net/http"

	"golang.org/x/crypto/bcrypt"
)

// Response is a generic response structure for our API endpoints.
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
	key, err := GenerateAPIKey(32)
	if err != nil {
		fmt.Println(err)
		http.Error(w, "Unable generate api keys..", http.StatusBadRequest)
		return
	}

	hashedKey, err := HashAPIKey(key)
	var secrets file.API_secrets
	session, err := Session.Store().Get(r, "session")
	secret := secrets.NewKey()
	secret.User = session.Values["user"].(string)
	if secret.User == "" {
		http.Error(w, "Unable generate api keys..", http.StatusForbidden)
		return
	}
	secret.Key = hashedKey
	secret.Name = r.URL.Query().Get("name")
	err = file.ReadJson(file.API_keys_file, &secrets)
	if err != nil {
		fmt.Println(err)
		http.Error(w, "Error getting existing keys..", http.StatusBadRequest)
		return
	}
	for _, k := range secrets.Keys {
		if k.Name == secret.Name {
			if err != nil {
				fmt.Println(err)
				http.Error(w, "Named key already exists!", http.StatusConflict)
				return
			}
		}
	}
	secrets.Keys = append(secrets.Keys, secret)
	err = file.WriteJson(file.API_keys_file, secrets)
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

	var secrets file.API_secrets
	err := file.ReadJson(file.API_keys_file, &secrets)
	if err != nil {
		http.Error(w, "Error retrieving API keys: "+err.Error(), http.StatusInternalServerError)
		return
	}

	type data struct {
		Name string `json:"name"` // The field should start with an uppercase letter
	}
	var apiList []data
	for _, k := range secrets.Keys {
		apiList = append(apiList, data{Name: k.Name})
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK) // Explicitly set status code for successful response
	json.NewEncoder(w).Encode(apiList)
}

// GenerateAPIKey creates a secure random API key
func GenerateAPIKey(length int) (string, error) {
	bytes := make([]byte, length)
	_, err := rand.Read(bytes)
	if err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}

func HashAPIKey(apiKey string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(apiKey), bcrypt.DefaultCost)
	return string(hash), err
}

func VerifyAPIKey(hashedKey, apiKey string) bool {
	return bcrypt.CompareHashAndPassword([]byte(hashedKey), []byte(apiKey)) == nil
}

func DeleteAPIKeyHandler(w http.ResponseWriter, r *http.Request) {
	if !Session.IsLoggedIn(w, r) {
		http.Redirect(w, r, "/login", http.StatusFound)
		return
	}

	var secrets file.API_secrets
	var new_secrets file.API_secrets

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

	err := file.ReadJson(file.API_keys_file, &secrets)
	if err != nil {
		fmt.Println(err)
		http.Error(w, "Error getting existing keys..", http.StatusBadRequest)
		return
	}

	session, err := Session.Store().Get(r, "session")
	username := session.Values["user"].(string)
	for _, k := range secrets.Keys {
		if k.Name == reqData.Data.Name && k.User == username {
			continue
		}
		new_secrets.Keys = append(new_secrets.Keys, k)
	}

	err = file.WriteJson(file.API_keys_file, new_secrets)
	if err != nil {
		fmt.Println(err)
		http.Error(w, "error saving new key..", http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	response := api_response{Status: true, Message: "Deleted api key.", Key: "nil"}
	json.NewEncoder(w).Encode(response)
}
