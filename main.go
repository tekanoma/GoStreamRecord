package main

import (
	"GoRecordurbate/modules/bot"
	"GoRecordurbate/modules/config"
	"fmt"
	"log"
	"os"
	"strings"
)

func init() {
	config.Path = "./config.json"
	config.C.Init()

}
func main() {
	args := os.Args
	if len(args) == 1 {
		fmt.Println("No args provided..")
		return
	}
	switch strings.ToLower(args[1]) {
	case "help": // To be added
	case "add":
		if len(args) < 3 {
			return
		}
		config.C.App.AddStreamer(args[2])
	case "del":
		if len(args) < 3 {
			return
		}
		config.C.App.RemoveStreamer(args[2])
	case "list":
		config.C.App.ListStreamers()
	case "import":
		if len(args) < 3 {
			return
		}
		config.C.App.ImportStreamers(args[2])
	case "export":
		config.C.App.ExportStreamers()
	case "start":
		logger := log.New(os.Stdout, "", log.LstdFlags)
		b := bot.NewBot(logger)
		b.Run()
	default:
		fmt.Println("Nothing to do..")
	}

}
