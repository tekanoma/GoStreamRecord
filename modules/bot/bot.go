package bot

import (
	"GoRecordurbate/modules/db"
	"GoRecordurbate/modules/file"
	"GoRecordurbate/modules/web/provider"
	"context"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"
	"sync"
	"time"
)

var Bot *bot

func Init() *bot {
	Bot = NewBot(log.New(os.Stdout, "lpg.log", log.LstdFlags))
	return Bot
}

// Bot encapsulates the recording bot’s state.
type bot struct {
	mux        sync.Mutex
	status     []BotStatus
	isFirstRun bool
	logger     *log.Logger
	// ctx is used to signal shutdown.
	ctx    context.Context
	cancel context.CancelFunc
}
type BotStatus struct {
	Enabled     bool
	StopStatus  bool
	WasRestart  bool
	Name        string    `json:"name"`
	Cmd         *exec.Cmd `json:"-"`
	IsRecording bool      `json:"isRecording"`
}

// NewBot creates a new Bot, sets up its cancellation context.
func NewBot(logger *log.Logger) *bot {
	ctx, cancel := context.WithCancel(context.Background())
	b := &bot{
		logger:     logger,
		ctx:        ctx,
		cancel:     cancel,
		status:     []BotStatus{},
		isFirstRun: true,
	}
	return b
}

// StopBot signals the bot to stop starting new recordings and then gracefully stops active processes.
func (b *bot) StopBot(streamerName string) {
	// Signal cancellation.
	b.cancel()
	log.Println("Stopping bot..")
	// Give current recorders time to finish (or exit gracefully).
	for i := range b.status {
		if b.status[i].Name == streamerName {
			b.status[i].StopStatus = true
		}
	}
}

// RecordLoop starts the main loop for a given streamer.
// It checks for online status, starts recording if not already recording, and listens for a shutdown signal.
func (b *bot) RecordLoop(streamerName string) {
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
			if !provider.Web.IsOnline(streamer.Name) {
				continue
			}
			b.mux.Lock()
			b.AddProcess(streamer.Name) // Assumes AddProcess safely appends to b.status.
			b.mux.Unlock()
			// Find the BotStatus for the streamer.
			for i2 := range b.status {
				if b.status[i2].Name != streamer.Name {
					continue
				}
				wg.Add(1)
				// Pass the index and streamer name into the closure to avoid capture issues.
				go func(status *BotStatus, sName string) {
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
							if !provider.Web.IsOnline(sName) {
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
					}
				}(&b.status[i2], streamer.Name)
			}
			if streamer.Name == streamerName {
				break
			}
		}
	}
	time.Sleep(time.Duration(db.Config.Settings.App.Loop_interval) * time.Second)
	wg.Wait()
}

// isRecorderActive returns true if a recorder for the given streamer is already running.
func (b *bot) isRecorderActive(streamerName string) bool {
	for _, rec := range b.status {
		if rec.Name == streamerName && rec.IsRecording {
			return true
		}
	}
	return false
}

// StopProcess sends a SIGINT to active recording processes and waits for them to finish.
func (b *bot) StopProcess(processName string) {
	b.mux.Lock()
	// Create a copy of status indices to avoid modification during iteration.
	statusCopy := make([]BotStatus, len(b.status))
	copy(statusCopy, b.status)
	b.mux.Unlock()

	for _, rec := range statusCopy {
		// Stop only the specified process (or all if processName is empty).
		if processName != "" && rec.Name != processName {
			continue
		}
		b.stopProcessIfRunning(rec)

	}
}

// writeYoutubeDLdb writes the youtube-dl configuration file.
func (b *bot) writeYoutubeDLdb() error {
	f, err := os.Create(file.YoutubeDL_configPath)
	if err != nil {
		return err
	}
	defer f.Close()

	folder := db.Config.Settings.App.Videos_folder
	dbLine := fmt.Sprintf("-o \"%s", folder) + "/%(id)s/%(title)s.%(ext)s\""
	_, err = f.Write([]byte(dbLine))
	return err
}

// startRecording starts a recording for the given streamer.
func (this *BotStatus) startRecording(streamerName string) {
	log.Printf("Starting recording for %s", streamerName)
	args := strings.Fields(db.Config.Settings.YoutubeDL.Binary)
	recordURL := fmt.Sprintf("https://chaturbate.com/%s/", streamerName)
	args = append(args, recordURL, "--config-location", file.YoutubeDL_configPath)
	this.Cmd = exec.Command(args[0], args[1:]...)

	// Start the recording process.
	if err := this.Cmd.Start(); err != nil {
		log.Printf("Error starting recording for %s: %v\n", streamerName, err)
	}
	this.Cmd.Wait()
}
