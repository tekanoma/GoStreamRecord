package login

import (
	"GoRecordurbate/modules/handlers/cookies"
	web_status "GoRecordurbate/modules/handlers/status"
	"encoding/json"
	"fmt"
	"net/http"

	"golang.org/x/crypto/bcrypt"
)

func PostLogin(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseMultipartForm(1024); err != nil {
		http.Error(w, "Invalid form submission", http.StatusBadRequest)
		return
	}
	username := r.FormValue("username")
	password := r.FormValue("password")
	resp := web_status.Response{}

	storedHash, ok := cookies.UserStore[username]
	if !ok {

		resp.Message = "Invalid credentials"
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(resp)
		return
	}
	if err := bcrypt.CompareHashAndPassword([]byte(storedHash), []byte(password)); err != nil {
		fmt.Println(err)
		resp.Message = "Invalid credentials"
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(resp)
		return
	}

	session, err := cookies.Session.Store().Get(r, "session")
	if err != nil {
		resp.Message = "Session error"
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(resp)
		return
	}
	session.Values["authenticated"] = true
	session.Values["user"] = username
	if err := session.Save(r, w); err != nil {
		resp.Message = "Could not save session"
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(resp)
		fmt.Println(err)
		//http.Error(w, "Could not save session", http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/", http.StatusFound)
}
