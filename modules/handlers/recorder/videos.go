package web_recorder

import (
	"GoRecordurbate/modules/config"
	"GoRecordurbate/modules/file"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
)

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

	err = filepath.Walk(config.Settings.App.Videos_folder, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() && file.IsVideoFile(info.Name()) {

			videos = append(videos, Video{URL: "/videos/" + filepath.Join(filepath.Base(filepath.Dir(path)), info.Name()), Name: info.Name()})
		}
		return nil
	})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if len(videos) == 0 {
		videos = append(videos, Video{URL: "", Name: "", NoVideos: fmt.Sprintf("No videos available. Try adding some to '%s'", config.Settings.App.Videos_folder)})

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
