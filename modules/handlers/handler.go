package handlers

import (
	"GoRecordurbate/modules/config"
	"GoRecordurbate/modules/file"
	web_config "GoRecordurbate/modules/handlers/config"
	"GoRecordurbate/modules/handlers/cookies"
	"GoRecordurbate/modules/handlers/login"
	web_recorder "GoRecordurbate/modules/handlers/recorder"
	web_status "GoRecordurbate/modules/handlers/status"
	"log"
	"net/http"
	"path/filepath"

	"golang.org/x/crypto/bcrypt"
)

func Handle() {
	// API endpoints

	http.Handle("/videos/", http.StripPrefix("/videos/", http.FileServer(http.Dir(config.C.App.Videos_folder))))

	http.HandleFunc("/api/add-streamer", web_config.AddStreamer)
	http.HandleFunc("/api/get-streamers", web_config.GetStreamers)
	http.HandleFunc("/api/remove-streamer", web_config.RemoveStreamer)
	http.HandleFunc("/api/control", web_recorder.ControlHandler)
	http.HandleFunc("/api/import", web_config.UploadHandler)
	http.HandleFunc("/api/export", web_config.DownloadHandler)
	http.HandleFunc("/api/status", web_status.StatusHandler)
	http.HandleFunc("/api/get-videos", web_recorder.GetVideos)

	password := "password"
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		log.Fatal("Error generating password hash:", err)
	}
	cookies.UserStore = map[string]string{
		"admin": string(hashedPassword),
	}

	fs := http.FileServer(http.Dir(filepath.Dir(file.Index_path)))

	http.HandleFunc("/login", func(w http.ResponseWriter, r *http.Request) {

		if r.Method == http.MethodGet {
			login.GetLogin(w, r)
		} else if r.Method == http.MethodPost {
			login.PostLogin(w, r)
		} else {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		session, err := cookies.Store.Get(r, "session")
		if err != nil {
			http.Error(w, "Session error", http.StatusInternalServerError)
			return
		}
		if auth, ok := session.Values["authenticated"].(bool); !ok || !auth {
			http.Redirect(w, r, "/login", http.StatusFound)
			//http.HandleFunc("/login", handlers.GetLogin)

		} else {
			fs.ServeHTTP(w, r)
		}
	})
}
