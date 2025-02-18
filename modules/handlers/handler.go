package handlers

import (
	"GoRecordurbate/modules/config"
	"GoRecordurbate/modules/file"
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"
	"time"
)

// Response is a generic response structure for our API endpoints.
type Response struct {
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

func StartWebUI() {
	Init()

	// API endpoints
	http.HandleFunc("/api/add-streamer", addStreamer)
	http.HandleFunc("/api/get-streamers", getStreamers)
	http.HandleFunc("/api/remove-streamer", removeStreamer)
	http.HandleFunc("/api/control", controlHandler)
	http.HandleFunc("/api/import", uploadHandler)
	http.HandleFunc("/api/export", downloadHandler)
	http.HandleFunc("/api/status", statusHandler)
	fs := http.FileServer(http.Dir(filepath.Dir(file.Index_path)))
	http.Handle("/", fs)

	server := &http.Server{
		Addr: fmt.Sprintf(":%d", config.C.App.Port),
	}

	// Channel to listen for termination signals
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)

	// Run the server in a separate goroutine
	go func() {
		log.Printf("Server starting on http://127.0.0.1:%d", config.C.App.Port)
		fmt.Printf("Server starting on http://127.0.0.1:%d\n", config.C.App.Port)
		if err := server.ListenAndServe(); err != http.ErrServerClosed {
			log.Fatalf("Server failed: %v", err)
		}
	}()

	// Wait for a termination signal
	<-stop
	log.Println("Shutting down server...")
	Bot.Stop()
	// Create a context with a timeout for graceful shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Fatalf("Server shutdown failed: %v", err)
	}

	log.Println("Server exited gracefully")
}
