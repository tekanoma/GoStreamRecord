package bot

import (
	"log"
	"strings"
	"sync"
	"time"
)

func (b *bot) Command(command string, name string) {
	if len(command) == 0 {
		log.Println("No command provided..")
		return
	}
	switch strings.ToLower(command) {
	case "start":
		for _, s := range b.ListRecorders() {
			if name == s.Name {
				log.Println("Bot already running..")
				return
			}
		}
		log.Println("Starting bot")
		b.mux.Lock()
		b.status.IsRunning = true
		b.isFirstRun = true
		b.mux.Unlock()
		b.RecordLoop(name)
	case "stop":
		is_running := false
		for _, s := range b.ListRecorders() {
			if name == s.Name {
				is_running = true

				break
			}
		}
		if !is_running && len(b.ListRecorders()) == 0 {
			log.Println("Bot is not running..")
			break
		}
		log.Println("Stopping bot")
		var wg sync.WaitGroup
		wg_was_added := false
		for _, s := range b.ListRecorders() {
			// Stop all

			if name == "" || s.Name == name {
				log.Println("Stopping:", name)
				wg_was_added = true
				wg.Add(1)
				go b.Stop(s.Name)
				// Stop single
			} else {
				log.Println("Not stopping..")
			}
		}
		if wg_was_added {

			wg.Wait()
		}
	case "restart":
		log.Println("Restarting bot")
		if b.status.IsRunning {
			b.Command("stop", name)
		}
		for len(b.status.Processes) != 0 {
			time.Sleep(1 * time.Second)
		}
		b.Command("start", name)
	default:
		log.Println("Nothing to do..")
	}

}
