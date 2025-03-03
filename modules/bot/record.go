package bot

import (
	"GoRecordurbate/modules/db"
	"GoRecordurbate/modules/file"
	"fmt"
	"log"
	"os/exec"
	"strings"
	"sync"
	"time"
)

// startRecording starts a recording for the given streamer.
func (this *Recorder) startRecording(streamerName string) {
	log.Printf("Starting recording for %s", streamerName)
	args := strings.Fields(db.Config.Settings.YoutubeDL.Binary)
	recordURL := fmt.Sprintf("%s%s/", this.Web.Url, streamerName)
	args = append(args, recordURL, "--config-location", file.YoutubeDL_configPath)
	this.Cmd = exec.Command(args[0], args[1:]...)

	// Start the recording process.
	if err := this.Cmd.Start(); err != nil {
		log.Printf("Error starting recording for %s: %v\n", streamerName, err)
	}
	this.Cmd.Wait()
}

// RecordLoop starts the main loop for a given streamer.
// It checks for online status, starts recording if not already recording, and listens for a shutdown signal.
func (b *controller) RecordLoop(streamerName string) {
	// Write youtube-dl db.
	if err := b.writeYoutubeDLdb(); err != nil {
		log.Println("Error writing youtube-dl db:", err)
		return
	}

	var wg sync.WaitGroup
	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()

	// Loop over configured streamers.
	for i1 := range db.Config.Streamers.Streamers {
		configIndex := i1
		streamer := db.Config.Streamers.Streamers[configIndex]
		if streamer.Name == streamerName || streamerName == "" {
			// Start a new recorder if one isn’t already running.
			if b.isRecorderActive(streamer.Name) {
				fmt.Println("Recorder was active")
				continue
			}
			b.mux.Lock()
			b.AddProcess(db.Config.Streamers.Streamers[i1].Provider, streamer.Name) 
			b.mux.Unlock()
			// Find the Recorder for the streamer.
			for i2 := range b.status {
				b.status[i2].Web.Site.Init(b.status[i2].WebType, b.status[i2].Name)
				//b.status[i2].Web.Site = provider.Init() //b.status[i2].Web.Type
				// Ensure correct name is being used.
				streamer.Name = b.status[i2].Web.Site.TrueName(streamer.Name)
				if b.status[i2].Name != streamer.Name {
					continue
				}
				wg.Add(1)
				// Pass the index and streamer name into the closure to avoid capture issues.
				go func(status *Recorder, sName string) {
					defer wg.Done()
					stopStatus := false
					for {
						// Exit the goroutine if the bot is cancelled.
						select {
						case <-b.ctx.Done():
							return
						default:
						}

						if stopStatus {

							b.StopProcess(sName)
							log.Println("Stopped!")
							// If not a restart, exit.
							b.mux.Lock()
							if !status.WasRestart {
								b.mux.Unlock()
								return
							}
							b.mux.Unlock()
							stopStatus = false
							status.WasRestart = false
						} else {
							b.mux.Lock()
							if b.isRecorderActive(sName) {
								b.mux.Unlock()
								return
							}
							b.mux.Unlock()
							// Optionally reload configuration.
							if db.Config.Settings.AutoReload {
								db.Config.Reload(file.Config_json_path, &db.Config.Settings)
							}
							log.Printf("Checking %s online status...", sName)
							if !status.Web.Site.IsOnline(sName) {
								log.Printf("Streamer %s is not online.", sName)
								return
							}
							log.Printf("Streamer %s is online!", sName)
							// Mark as recording.
							b.mux.Lock()
							status.IsRecording = true
							b.mux.Unlock()

							status.startRecording(sName)

							b.mux.Lock()
							stopStatus = status.StopStatus
							status.IsRecording = false
							status.StopStatus = true
							b.mux.Unlock()

							log.Printf("Recording for %s finished", sName)
							stopStatus = true
						}
						time.Sleep(time.Duration(db.Config.Settings.App.Loop_interval) * time.Minute)
					}
				}(&b.status[i2], streamer.Name)
			}
			if streamer.Name == streamerName {
				break
			}
		}
	}
	time.Sleep(time.Duration(db.Config.Settings.App.Loop_interval) * time.Minute)
	wg.Wait()
}
