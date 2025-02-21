package bot

import (
	"GoRecordurbate/modules/config"
	"fmt"
	"log"
	"strings"
	"time"
)

func (b *bot) Command(command string) {
	if len(command) == 0 {
		fmt.Println("No command provided..")
		return
	}
	switch strings.ToLower(command) {
	case "import":
		log.Println("Importing streamers..")
		config.C.App.ImportStreamers(config.C.App.ExportPath)
	case "export":
		log.Println("Exporting streamers..")
		config.C.App.ExportStreamers()
	case "start":
		log.Println("Starting bot")
		go b.Run()
	case "stop":
		log.Println("Stopping bot")
		b.Stop()
	case "restart":
		log.Println("Restarting bot")
		b.Stop()
		time.Sleep(1 * time.Second)
		go b.Run()
	case "start_monitoring":
		log.Println("Monitoring not implemented")
	case "restarting bot":
		log.Println("Stopping bot")
		b.ResetBot()
		go b.Run()
	default:
		fmt.Println("Nothing to do..")
	}

}
