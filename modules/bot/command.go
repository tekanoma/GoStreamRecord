package bot

import (
	"context"
	"log"
	"strings"
	"sync"
	"syscall"
)

func (b *bot) Command(command string, name string) {
	if len(command) == 0 {
		log.Println("No command provided..")
		return
	}
	switch strings.ToLower(command) {
	case "start":
		// If the bot was previously stopped, reinitialize the context.
		if b.ctx.Err() != nil {
			b.ctx, b.cancel = context.WithCancel(context.Background())
		}

		for i, s := range b.status {
			if name == s.Name && s.Cmd != nil {
				log.Println("Bot already running..")
				return
			} else if s.Cmd == nil {
				if err := s.Cmd.Process.Signal(syscall.Signal(0)); err != nil {
					log.Printf("Process for %s has stopped", s.Name)
					b.status = append(b.status[:i], b.status[i+1:]...)
					i--
				}
			}
		}
		log.Println("Starting bot")
		go b.RecordLoop(name)
	case "stop":
		is_running := len(b.status) != 0

		if !is_running && len(b.status) == 0 {
			log.Println("[bot] Stopped recording for")
			break
		}
		log.Println("Stopping bot")
		var wg sync.WaitGroup
		// Iterate over a copy of the status slice to avoid closure capture issues.
		for i, s := range b.status {
			// Stop only the specified process (or all if name is empty).
			if name == "" || s.Name == name {
				b.status[i].WasRestart = true
				b.stopProcessIfRunning(b.status[i])
				sName := s.Name
				wg.Add(1)
				go func(n string) {
					defer wg.Done()
					b.StopProcess(n)
				}(sName)
			} else {
				log.Println("Not stopping..")
			}
		}
		wg.Wait()

		b.checkProcesses()
	case "restart":
		log.Println("Restarting bot")
		recorders := []string{}
		// Before restarting, reinitialize the context so RecordLoop doesn't exit immediately.
		b.ctx, b.cancel = context.WithCancel(context.Background())

		if name != "" {
			// Stop a single process.
			process := getProcess(name, b)

			b.Command("stop", process.Name)
			recorders = append(recorders, name)

		} else {
			var wg sync.WaitGroup
			// Stop all running recorders.
			// Create a copy of b.status to avoid data races when stopping processes.
			b.mux.Lock()
			statusCopy := make([]BotStatus, len(b.status))
			copy(statusCopy, b.status)
			b.mux.Unlock()
			for _, s := range statusCopy {
				b.mux.Lock()

				// Mark that the process is being restarted.
				// (Assuming b.status is the source of truth; you might also update the copy)
				for i, rec := range b.status {
					if rec.Name == s.Name {
						b.status[i].WasRestart = true
						b.stopProcessIfRunning(b.status[i])
						break
					}
				}
				b.mux.Unlock()
				wg.Add(1)
				recorders = append(recorders, s.Name)
				go func(n string) {
					b.Command("stop", n)
					log.Println("Stopped", n)
					wg.Done()
				}(s.Name)
			}
			wg.Wait()
		}

		// Start all recorders that were stopped.
		for _, recName := range recorders {
			go b.RecordLoop(recName)
		}
	default:
		log.Println("Nothing to do..")
	}
}

func (b *bot) stopProcessIfRunning(bs BotStatus) {

	for i, s := range b.status {
		if bs.Cmd != nil && s.Name == bs.Name {
			b.status[i].StopStatus = true
			if err := s.Cmd.Process.Signal(syscall.Signal(0)); err != nil {
				i--
			}
			b.status = append(b.status[:i], b.status[i+1:]...)
			break
		}
		if s.Cmd == nil && s.Name == bs.Name {
			b.status[i].StopStatus = true
			b.status = append(b.status[:i], b.status[i+1:]...)
			break
		}
	}

	log.Printf("Process for %s has stopped", bs.Name)
}
