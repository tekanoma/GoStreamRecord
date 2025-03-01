package main

import (
	"GoRecordurbate/modules/bot"
	"GoRecordurbate/modules/db"
	"GoRecordurbate/modules/file"
	"GoRecordurbate/modules/handlers"
	"GoRecordurbate/modules/handlers/cookies"
	"GoRecordurbate/modules/handlers/login"
	"context"
	_ "embed"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/joho/godotenv"
)

// Embed static HTML files
//
//go:embed internal/web/index.html
var IndexHTML string

//go:embed internal/web/login.html
var LoginHTML string

func init() {
	handlers.IndexHTML = IndexHTML
	handlers.LoginHTML = LoginHTML
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatalf("Some error occured. Err: %s", err)
	}
	os.Mkdir("./output", 0755)
	cookies.Session = cookies.New([]byte(os.Getenv("SESSION_KEY")))
	file.InitLog(file.Log_path)
	bot.Init()

}

func main() {
	if len(os.Args) > 1 {
		if os.Args[1] != "reset-pwd" {
			fmt.Println("Usage: ./GoRecordurbate reset-pwd <username> <new-password>")
			fmt.Println("Otherwise run the server without any arguments.")
			return
		}

		if len(os.Args) <= 2 {
			fmt.Println("No username provided.")
			fmt.Println("Usage: ./GoRecordurbate reset-pwd <username> <new-password>")
			fmt.Println("Otherwise run the server without any arguments.")
			return
		}

		username := os.Args[2]
		if len(os.Args) <= 3 {
			fmt.Println("No new password provided.")
			fmt.Println("Usage: ./GoRecordurbate reset-pwd <username> <new-password>")
			fmt.Println("Otherwise run the server without any arguments.")
			return
		}

		newPassword := os.Args[3]
		userFound := false

		for i, u := range db.Users.Users {
			if u.Name == username {
				db.Users.Users[i].Key = string(login.HashedPassword(newPassword))
				userFound = true
				break
			}
		}

		if !userFound {
			log.Println("No matching username found.")
			fmt.Println("No matching username found.")
			return
		}
		db.Update(file.Users_json_path, db.Users)
		log.Println("Password updated for", username)
		fmt.Println("Password updated for", username)
		return // Exit after resetting password

	}

	//http.Handle("/", fs)
	handlers.Handle()
	server := &http.Server{
		Addr: fmt.Sprintf(":%d", db.Settings.App.Port),
	}

	// Channel to listen for termination signals
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)

	// Run the server in a separate goroutine
	go func() {
		log.Printf("Server starting on http://127.0.0.1:%d", db.Settings.App.Port)
		fmt.Printf("Server starting on http://127.0.0.1:%d\n", db.Settings.App.Port)
		if err := server.ListenAndServe(); err != http.ErrServerClosed {
			log.Fatalf("Server failed: %v", err)
		}
	}()

	// Wait for a termination signal
	<-stop
	log.Println("Shutting down server...")
	bot.Bot.StopBot("")
	// Create a context with a timeout for graceful shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Fatalf("Server shutdown failed: %v", err)
	}

	log.Println("Server exited gracefully")
}
