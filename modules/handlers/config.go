package handlers

import (
	"GoRecordurbate/modules/config"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"path/filepath"
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
	log.Printf("Uploaded File: %s, Size: %d bytes\n", handler.Filename, len(fileContent))

	resp := Response{
		Message: fmt.Sprintf("File '%s' uploaded successfully", handler.Filename),
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
