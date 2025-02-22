package handlers

import (
	"GoRecordurbate/modules/config"
	"GoRecordurbate/modules/file"
	web_config "GoRecordurbate/modules/handlers/config"
	"GoRecordurbate/modules/handlers/cookies"
	"GoRecordurbate/modules/handlers/login"
	web_recorder "GoRecordurbate/modules/handlers/recorder"
	web_status "GoRecordurbate/modules/handlers/status"
	"net/http"
	"path/filepath"
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
	if cookies.UserStore == nil {
		cookies.UserStore = make(map[string]string)
	}

	for _, u := range config.Users.Users {
		cookies.UserStore[u.Name] = u.Key
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
		session, err := cookies.Session.Store().Get(r, "session")
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
