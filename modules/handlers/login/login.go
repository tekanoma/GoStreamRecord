package login

import (
	"GoRecordurbate/modules/file"
	"GoRecordurbate/modules/handlers/cookies"
	web_response "GoRecordurbate/modules/handlers/response"
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"
	"os"

	"golang.org/x/crypto/bcrypt"
)

var tmpl *template.Template

func init() {

	contentBytes, _ := os.ReadFile(file.Login_path)

	tmpl = template.Must(template.New("login").Parse(string(contentBytes)))
}

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
	resp := web_response.Response{}

	storedHash, ok := cookies.UserStore[username]
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

	session, err := cookies.Store.Get(r, "session")
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
