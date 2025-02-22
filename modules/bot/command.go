package bot

import (
	"context"
	"fmt"
	"log"
	"strings"
)

func (b *bot) Command(command string) {
	b.ctx, b.cancel = context.WithCancel(context.Background())
	if len(command) == 0 {
		fmt.Println("No command provided..")
		return
	}
	switch strings.ToLower(command) {
	case "start":
		if b.isRunning {
			log.Println("Bot already running..")
			break
		}
		log.Println("Starting bot")
		go b.RecordLoop()
	case "stop":
		if !b.isRunning {
			log.Println("Bot already running..")
			break
		}
		log.Println("Stopping bot")

		b.Stop()
	case "restart":
		log.Println("Restarting bot")
		if b.isRunning {
			b.Command("stop")
		}
		b.Command("start")
	default:
		fmt.Println("Nothing to do..")
	}

}
