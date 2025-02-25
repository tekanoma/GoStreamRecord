package bot

import (
	"log"
	"strings"
	"time"
)

func (b *bot) Command(command string) {
	if len(command) == 0 {
		log.Println("No command provided..")
		return
	}
	switch strings.ToLower(command) {
	case "start":
		if b.status.IsRunning {
			log.Println("Bot already running..")
			break
		}
		log.Println("Starting bot")
		b.mux.Lock()
		b.status.IsRunning = true
		b.isFirstRun = true
		b.mux.Unlock()
		go b.RecordLoop("")
	case "stop":
		if !b.status.IsRunning {
			log.Println("Bot is not running..")
			break
		}
		log.Println("Stopping bot")
		b.stopRecording <- true
	case "restart":
		log.Println("Restarting bot")
		if b.status.IsRunning {
			b.Command("stop")
		}
		for len(b.status.Processes) != 0 {
			time.Sleep(1 * time.Second)
		}
		b.Command("start")
	default:
		log.Println("Nothing to do..")
	}

}
