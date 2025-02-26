package bot

import (
	"GoRecordurbate/modules/config"
	"GoRecordurbate/modules/file"
	"fmt"
	"log"
	"os/exec"
	"strings"
)

// runRecordLoop starts a recording for the given streamer (if online) and waits for the process to finish.
func (b *bot) runRecordLoop(streamerName string) {


	log.Printf("[bot]: Checking %s room status...", streamerName)
	if !b.IsOnline(streamerName) {
		log.Printf("[bot]: Streamer %s is not online.", streamerName)
		return
	}

	log.Printf("[bot]: Starting recording for %s", streamerName)
	args := strings.Fields(config.Settings.YoutubeDL.Binary)
	recordURL := fmt.Sprintf("https://chaturbate.com/%s/", streamerName)
	args = append(args, recordURL, "--config-location", file.YoutubeDL_configPath)
	cmd := exec.Command(args[0], args[1:]...)

	// Start the recording process.
	if err := cmd.Start(); err != nil {
		log.Printf("[bot]: Error starting recording for %s: %v", streamerName, err)
		return
	}

	// Add the process to our list.
	b.mux.Lock()
	b.status.Processes = append(b.status.Processes, StreamerStatus{
		Name:        streamerName,
		IsRecording: true,
		Cmd:         cmd,
	})
	b.mux.Unlock()

	// Wait for the command to finish.
	cmd.Wait()
	log.Printf("[bot]: Recording for %s finished", streamerName)

	// Remove this process from our list.
	b.mux.Lock()
	for i, p := range b.status.Processes {
		if p.Name == streamerName && p.Cmd == cmd {
			b.status.Processes = append(b.status.Processes[:i], b.status.Processes[i+1:]...)
			break
		}
	}
	b.mux.Unlock()
}
