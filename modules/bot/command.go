package bot

import (
	"GoRecordurbate/modules/config"
	"fmt"
	"strings"
	"time"
)

func (b *Bot) Command(command string) {
	if len(command) == 0 {
		fmt.Println("No command provided..")
		return
	}
	switch strings.ToLower(command) {
	case "help": // To be added
	case "add":
		if len(command) < 3 {
			return
		}
		config.C.App.AddStreamer(command)
	case "del":
		if len(command) < 3 {
			return
		}
		config.C.App.RemoveStreamer(command)
	case "list":
		config.C.App.ListStreamers()
	case "import":
		if len(command) < 3 {
			return
		}
		config.C.App.ImportStreamers(command)
	case "export":
		config.C.App.ExportStreamers()
	case "start":
		fmt.Println("Starting", b)
		b.running = true
		go b.Run()
	case "stop":
		fmt.Println("Stopping", b)
		b.Stop()
	case "restart":
		fmt.Println("Restarting", b)
		b.Stop()
		time.Sleep(1 * time.Second)
		b.running = true
		go b.Run()
	default:
		fmt.Println("Nothing to do..")
	}

}
