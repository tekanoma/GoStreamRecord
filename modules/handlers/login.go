package handlers

import (
	"GoRecordurbate/modules/file"
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"
	"os"

	"github.com/gorilla/sessions"
	"golang.org/x/crypto/bcrypt"
)

var Store = sessions.NewCookieStore([]byte("very-secret-key-here"))
var tmpl *template.Template

func init() {

	contentBytes, _ := os.ReadFile(file.Login_path)

	tmpl = template.Must(template.New("login").Parse(string(contentBytes)))

	Store.Options = &sessions.Options{
		Path:     "/",
		HttpOnly: true,
		Secure:   false, // Set to false for HTTP
		MaxAge:   3600,
	}
}

var UserStore map[string]string

func GetLogin(w http.ResponseWriter, r *http.Request) {

	if err := tmpl.Execute(w, nil); err != nil {
		http.Error(w, "Server error", http.StatusInternalServerError)
	}

}

func PostLogin(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseMultipartForm(1024); err != nil {
		http.Error(w, "Invalid form submission", http.StatusBadRequest)
		return
	}
	username := r.FormValue("username")
	password := r.FormValue("password")
	resp := Response{}

	storedHash, ok := UserStore[username]
	if !ok {

		fmt.Println(username, password)
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

	session, err := Store.Get(r, "session")
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
		//http.Error(w, "Could not save session", http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/", http.StatusFound)
}
