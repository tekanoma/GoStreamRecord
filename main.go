package main

import (
	"GoRecordurbate/modules/bot"
	"GoRecordurbate/modules/config"
	"GoRecordurbate/modules/file"
	"GoRecordurbate/modules/handlers"
	"GoRecordurbate/modules/handlers/cookies"
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/joho/godotenv"
)

func init() {

	err := godotenv.Load(".env")
	if err != nil {
		log.Fatalf("Some error occured. Err: %s", err)
	}
	cookies.Session = cookies.New([]byte(os.Getenv("SESSION_KEY")))
	file.InitLog(file.Log_path)
	bot.Init()

}

func main() {
	//http.Handle("/", fs)
	handlers.Handle()
	server := &http.Server{
		Addr: fmt.Sprintf(":%d", config.Settings.App.Port),
	}

	// Channel to listen for termination signals
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)

	// Run the server in a separate goroutine
	go func() {
		log.Printf("Server starting on http://127.0.0.1:%d", config.Settings.App.Port)
		fmt.Printf("Server starting on http://127.0.0.1:%d\n", config.Settings.App.Port)
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
