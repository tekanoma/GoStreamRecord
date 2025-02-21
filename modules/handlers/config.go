package handlers

import (
	"GoRecordurbate/modules/config"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"path/filepath"
	"strings"
)

// uploadHandler handles POST /api/upload.
// It reads an uploaded file and returns a dummy success response.
func UploadHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Only POST allowed", http.StatusMethodNotAllowed)
		return
	}

	// Limit the size of the incoming request to 10MB
	if err := r.ParseMultipartForm(10 << 20); err != nil {
		http.Error(w, "Unable to parse form", http.StatusBadRequest)
		return
	}

	// Retrieve file from posted form-data
	file, handler, err := r.FormFile("file")
	if err != nil {
		http.Error(w, "Error retrieving the file", http.StatusBadRequest)
		return
	}
	defer file.Close()

	if filepath.Ext(handler.Filename) != ".txt" {
		return
	}

	// For demonstration, we'll read the file's contents (but not store it)
	fileContent, err := io.ReadAll(file)
	if err != nil {
		http.Error(w, "Error reading file", http.StatusInternalServerError)
		return
	}
	newStreamers := []string{}
	for _, line := range strings.Split(string(fileContent), "\n") {
		exist := false
		for _, s := range config.C.App.Streamers {
			if line == s.Name {
				exist = true
			}
		}
		if exist {
			continue
		}
		newStreamers = append(newStreamers, line)
	}
	for _, line := range newStreamers {
		if len(line) == 0 {
			continue
		}
		config.C.App.AddStreamer(line)
	}
	resp := Response{
		Message: fmt.Sprintf("Added %d new streamers!", len(newStreamers)),
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

// downloadHandler handles GET /api/download.
// It sends a dummy file for download.
func DownloadHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Only GET allowed", http.StatusMethodNotAllowed)
		return
	}

	// Dummy file content
	fileContent := ""

	for _, s := range config.C.App.Streamers {
		fileContent = fileContent + s.Name + "\n"
	}
	w.Header().Set("Content-Type", "application/octet-stream")
	w.Header().Set("Content-Disposition", "attachment; filename=export.txt")
	w.Write([]byte(fileContent))
}
