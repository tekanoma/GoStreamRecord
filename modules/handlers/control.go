package handlers

import (
	"GoRecordurbate/modules/bot"
	"GoRecordurbate/modules/config"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

// dcodes a JSON payload with a "command" field (start, stop, or restart)
// and returns a dummy response.
func ControlHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Only POST allowed", http.StatusMethodNotAllowed)
		return
	}

	type RequestData struct {
		Command string `json:"command"`
	}
	var reqData RequestData
	if err := json.NewDecoder(r.Body).Decode(&reqData); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}
	
	bot.Bot.Command(reqData.Command)
	resp := Response{
		Message: fmt.Sprintf("Control command '%s' executed", reqData.Command),
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)

}

type Video struct {
	URL      string `json:"url"`
	Name     string `json:"name"`
	NoVideos string `json:"error"`
}

func GetVideos(w http.ResponseWriter, r *http.Request) {
	videos := []Video{}

	page, err := strconv.Atoi(r.URL.Query().Get("page"))
	if err != nil || page < 1 {
		page = 1
	}

	err = filepath.Walk(config.C.App.Videos_folder, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() && isVideoFile(info.Name()) {

			videos = append(videos, Video{URL: "/videos/" + filepath.Join(filepath.Base(filepath.Dir(path)), info.Name()), Name: info.Name()})
		}
		return nil
	})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if len(videos) == 0 {
		videos = append(videos, Video{URL: "", Name: "", NoVideos: fmt.Sprintf("No videos available. Try adding some to '%s'", config.C.App.Videos_folder)})

	}
	pageSize := 10
	start := (page - 1) * pageSize
	end := start + pageSize
	if start >= len(videos) {
		start = len(videos)
	}
	if end > len(videos) {
		end = len(videos)
	}

	paginatedVideos := videos[start:end]

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(paginatedVideos)
}

// isVideoFile returns true if the file extension indicates a video file.
func isVideoFile(filename string) bool {
	extensions := []string{".mp4", ".avi", ".mov", ".mkv"}
	lower := strings.ToLower(filename)
	for _, ext := range extensions {
		if strings.HasSuffix(lower, ext) {
			return true
		}
	}
	return false
}
