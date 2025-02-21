package main

import (
	"GoRecordurbate/modules/bot"
	"GoRecordurbate/modules/config"
	"GoRecordurbate/modules/file"
	"GoRecordurbate/modules/handlers"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func init() {
	_, err := os.Stat(file.Config_path)
	if os.IsNotExist(err) {
		fmt.Println("No config found. Generating.. Please fill in details in:", file.Config_path)
		f, _ := os.Create(file.Config_path)
		tmp := config.Config{}
		b, _ := json.Marshal(&tmp)
		f.Write(b)
		os.Exit(0)
	}
	file.InitLog(file.Log_path)
	bot.Init()
	config.C.Init()

}

func main() {
	//http.Handle("/", fs)
	handlers.Handle()
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
	bot.Bot.Stop()
	// Create a context with a timeout for graceful shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Fatalf("Server shutdown failed: %v", err)
	}

	log.Println("Server exited gracefully")
}
