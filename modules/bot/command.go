package bot

import (
	"GoRecordurbate/modules/config"
	"fmt"
	"log"
	"strings"
	"time"
)

func (b *Bot) Command(command string) {
	if len(command) == 0 {
		fmt.Println("No command provided..")
		return
	}
	switch strings.ToLower(command) {
	case "add":
		log.Println("Adding ")
		config.C.App.AddStreamer("STREAMER_NAME_PLACEHOLDER")
	case "del":
		log.Println("Removing ")
		config.C.App.RemoveStreamer("STREAMER_NAME_PLACEHOLDER")
	case "import":
		log.Println("Importing streamers..")
		config.C.App.ImportStreamers(config.C.App.ExportPath)
	case "export":
		log.Println("Exporting streamers..")
		config.C.App.ExportStreamers()
	case "start":
		log.Println("Starting bot")
		if b.running {
			log.Println("Bot already started. Use 'Restart'.")
			return
		}
		b.running = true
		go b.Run()
	case "stop":
		log.Println("Stopping bot")
		go b.Stop()
	case "restart":
		log.Println("Restarting bot")
		b.Stop()
		time.Sleep(1 * time.Second)
		b.running = true
		go b.Run()
	default:
		fmt.Println("Nothing to do..")
	}

}
