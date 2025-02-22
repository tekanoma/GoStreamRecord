package handlers

import (
	"GoRecordurbate/modules/config"
	"GoRecordurbate/modules/file"
	web_config "GoRecordurbate/modules/handlers/config"
	"GoRecordurbate/modules/handlers/cookies"
	"GoRecordurbate/modules/handlers/login"
	web_recorder "GoRecordurbate/modules/handlers/recorder"
	web_status "GoRecordurbate/modules/handlers/status"
	"fmt"
	"log"
	"net/http"
	"os"
	"text/template"
)

func Handle() {
	// API endpoints

	http.Handle("/videos/", http.StripPrefix("/videos/", http.FileServer(http.Dir(config.Settings.App.Videos_folder))))

	http.HandleFunc("/api/add-streamer", web_config.AddStreamer)
	http.HandleFunc("/api/get-streamers", web_config.GetStreamers)
	http.HandleFunc("/api/remove-streamer", web_config.RemoveStreamer)
	http.HandleFunc("/api/control", web_recorder.ControlHandler)
	http.HandleFunc("/api/import", web_config.UploadHandler)
	http.HandleFunc("/api/export", web_config.DownloadHandler)
	http.HandleFunc("/api/status", web_status.StatusHandler)
	http.HandleFunc("/api/get-videos", web_recorder.GetVideos)
	http.HandleFunc("/api/logs", web_recorder.HandleLogs)

	http.HandleFunc("/api/get-users", web_config.GetUsers)
	http.HandleFunc("/api/add-user", web_config.AddUser)
	http.HandleFunc("/api/update-user", web_config.UpdateUsers)
	fmt.Println(len(os.Args))
	if len(os.Args) > 1 {
		if os.Args[1] != "reset-pwd" {
			log.Println("Nothing to do..")
			fmt.Println("Nothing to do..")
			return
		}

		if len(os.Args) <= 2 {
			log.Println("No username provided.")
			fmt.Println("No username provided.")
			return
		}

		username := os.Args[2]
		if len(os.Args) <= 3 {
			log.Println("No new password provided.")
			fmt.Println("No new password provided.")
			return
		}

		newPassword := os.Args[3]
		userFound := false

		for i, u := range config.Users.Users {
			if u.Name == username {
				config.Users.Users[i].Key = string(login.HashedPassword(newPassword))
				userFound = true
				break
			}
		}

		if !userFound {
			log.Println("No matching username found.")
			fmt.Println("No matching username found.")
		}
		config.Update(file.Users_json_path, config.Users)
		log.Println("Password updated!")
		fmt.Println("Password updated!")
		os.Exit(0)

	}

	if cookies.UserStore == nil {
		cookies.UserStore = make(map[string]string)
	}

	for _, u := range config.Users.Users {
		cookies.UserStore[u.Name] = u.Key
	}
	//	fs := http.FileServer(http.Dir(filepath.Dir(file.Index_path)))

	http.HandleFunc("/login", func(w http.ResponseWriter, r *http.Request) {

		if r.Method == http.MethodGet {
			GetLogin(w, r)
		} else if r.Method == http.MethodPost {
			login.PostLogin(w, r)
		} else {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		session, err := cookies.Session.Store().Get(r, "session")
		if err != nil {
			http.Error(w, "Session error", http.StatusInternalServerError)
			return
		}
		if auth, ok := session.Values["authenticated"].(bool); !ok || !auth {
			http.Redirect(w, r, "/login", http.StatusFound)
			//http.HandleFunc("/login", handlers.GetLogin)

		} else {
			GetIndex(w, r)
		}
	})
}

var IndexHTML, LoginHTML string

type Template struct {
	W    http.ResponseWriter
	Tmpl *template.Template
}

func (t *Template) Execute(data any) error {
	return t.Tmpl.Execute(t.W, data)
}

func GetIndex(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.New("index").Parse(IndexHTML))
	indexTemplate := Template{W: w, Tmpl: tmpl}
	if err := indexTemplate.Execute(nil); err != nil {
		http.Error(w, "Server error", http.StatusInternalServerError)
	}
}

func GetLogin(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.New("login").Parse(LoginHTML))
	loginTemplate := Template{W: w, Tmpl: tmpl}
	if err := loginTemplate.Execute(nil); err != nil {
		http.Error(w, "Server error", http.StatusInternalServerError)
	}
}
