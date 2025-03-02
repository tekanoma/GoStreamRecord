package bot

import (
	"log"
	"syscall"
)

func getProcesss(name string, b *bot) BotStatus {
	b.mux.Lock()
	defer b.mux.Unlock()
	for _, s := range b.status {
		if name == s.Name {
			return s
		}
	}
	return BotStatus{StopStatus: false, Name: name, IsRecording: false, Cmd: nil}
}

func (b *bot) AddProcess(streamerName string) {
	// Only add if not already present
	for _, rec := range b.status {
		if rec.Name == streamerName {
			return
		}
	}
	b.status = append(b.status, BotStatus{Name: streamerName})
}

func getProcess(name string, b *bot) BotStatus {
	b.mux.Lock()
	defer b.mux.Unlock()
	for _, s := range b.status {
		if name == s.Name {
			return s
		}
	}
	return BotStatus{StopStatus: false, Name: name, IsRecording: false, Cmd: nil}
}

func (b *bot) Status(name string) BotStatus {
	return getProcess(name, b)
}

// ListRecorders returns the current list of recorder statuses.
func (b *bot) ListRecorders() []BotStatus {
	b.mux.Lock()
	defer b.mux.Unlock()
	return b.status
}
func (b *bot) StopRunningEmpty() {
	b.checkProcesses()
}

// checkProcesses looks through the list of processes and removes any that have finished.
func (b *bot) checkProcesses() int {
	b.mux.Lock()
	defer b.mux.Unlock()
	for i := 0; i < len(b.status); i++ {
		rec := b.status[i]
		// Use signal 0 to check if process is still running.
		if !rec.StopStatus {
			continue
		}
		if err := rec.Cmd.Process.Signal(syscall.Signal(0)); err != nil {
			log.Printf("Process for %s has stopped", rec.Name)
			b.status = append(b.status[:i], b.status[i+1:]...)
			i--
		}
	}
	return len(b.status)
}
