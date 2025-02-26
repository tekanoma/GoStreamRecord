package bot

import (
	"GoRecordurbate/modules/config"
	"GoRecordurbate/modules/file"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"os/exec"
	"os/signal"
	"strings"
	"sync"
	"syscall"
	"time"
)

var Bot *bot

func Init() {
	Bot = NewBot(log.New(os.Stdout, "lpg.log", log.LstdFlags))

}

// Interface can be expanded as needed.
type Interface interface {
	AppendStreamer(streamer string)
}

// Bot encapsulates the recording bot’s state.
type bot struct {
	mux        sync.Mutex
	status     BotStatus
	isFirstRun bool
	logger     *log.Logger
	Interface
	stopProcessName chan string
	// ctx is used to signal shutdown.
	ctx    context.Context
	cancel context.CancelFunc
}
type BotStatus struct {
	IsRunning bool             `json:"isRunning"`
	Processes []StreamerStatus `json:"processes"`
}

type StreamerStatus struct {
	Name        string    `json:"name"`
	IsRecording bool      `json:"isRecording"`
	Cmd         *exec.Cmd `json:"-"`
}

// NewBot creates a new Bot, sets up its cancellation context, and registers a signal handler.
func NewBot(logger *log.Logger) *bot {

	ctx, cancel := context.WithCancel(context.Background())
	s := BotStatus{
		IsRunning: false,
	}
	b := &bot{
		stopProcessName: make(chan string),
		logger:          logger,
		ctx:             ctx,
		cancel:          cancel,
		status:          s,
		isFirstRun:      true,
	}

	// Register to catch SIGINT and SIGTERM and trigger Stop.
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		s := <-sigs
		logger.Printf("Caught signal %v, stopping", s)
		b.Stop("")
	}()
	return b
}

func (b *bot) Status() BotStatus {
	return b.status
}
func (b *bot) Stop(processName string) {
	b.stopProcessName <- processName
}

func (b *bot) AppendStreamer(name string) {
	b.mux.Lock()
	defer b.mux.Unlock()
	b.status.Processes = append(b.status.Processes, StreamerStatus{Name: name, IsRecording: false})
}

// ListRecorders returns the current list of recorder statuses.
func (b *bot) ListRecorders() []StreamerStatus {
	b.mux.Lock()
	defer b.mux.Unlock()
	return b.status.Processes
}

// StopProcess stops all the recordings if no streamer name is provided.
func (b *bot) StopProcess(streamerName string) {
	// Signal cancellation.
	if streamerName == "" {
		log.Println("Stopping all recordings")
	} else {
		log.Println("Stopping recording for", streamerName)
	}
	// Give current recorders time to finish (or exit gracefully).
	b.stopActiveProcesses(streamerName)
}

// Stop signals the bot to stop starting new recordings and then gracefully stops active processes.
//
// streamerName can be used to stop a single recording.
func (b *bot) StopBot(streamerName string) {
	// Signal cancellation.
	b.cancel()
	log.Println("Stopping bot..")
	// Give current recorders time to finish (or exit gracefully).
	b.stopActiveProcesses("")
}

// IsRoomPublic checks if a given room is public by sending a POST request.
func (b *bot) IsRoomPublic(username string) bool {
	// Wait for the configured rate limit.
	time.Sleep(time.Duration(config.Settings.App.RateLimit.Time) * time.Second)
	urlStr := "https://chaturbate.com/get_edge_hls_url_ajax/"
	data := url.Values{}
	data.Set("room_slug", username)
	req, err := http.NewRequest("POST", urlStr, strings.NewReader(data.Encode()))
	if err != nil {
		log.Printf("Error creating request: %v", err)
		return false
	}
	req.Header.Add("X-Requested-With", "XMLHttpRequest")
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Printf("Error making the request: %v", err)
		return false
	}
	defer resp.Body.Close()

	var res struct {
		Success    bool   `json:"success"`
		RoomStatus string `json:"room_status"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&res); err != nil {
		log.Printf("Error decoding JSON: %v", err)
		return false
	}
	return res.Success && res.RoomStatus == "public"
}

// IsOnline checks if the streamer is online by checking if a thumbnail is available from the stream.
func (b *bot) IsOnline(username string) bool {
	// Short delay before making the call.

	//Check once if a thumbnail is available
	urlStr := "https://jpeg.live.mmcdn.com/stream?room=" + username
	resp, err := http.Get(urlStr)
	if err != nil {
		log.Printf("Error in making request: %v", err)
		return false
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK { // Streamer is not online if response if not 200
		return false

	}

	return true
}

// Run starts the main loop of the Bot.
// Reloads the configuration, checks running processesn for each streamer, starts a recorder if one isn’t already running.
// It also checks for a stop signal and waits for active recordings to finish before stopping.
// Streamer name is optional and can be used to start a single recorder.
func (b *bot) RecordLoop(streamerName string) {
	b.status.IsRunning = true
	b.isFirstRun = true
	is_single_run := false
	if streamerName != "" {
		is_single_run = true
	}
	// Write youtube-dl config.
	if err := b.writeYoutubeDLConfig(); err != nil {
		log.Println("Error writing youtube-dl config:", err)
		return
	}

	var wg sync.WaitGroup
	ticker := time.NewTicker(time.Duration(1) * time.Second)

	defer ticker.Stop()

	// Main loop.
	for {
		select {
		case name := <-b.stopProcessName:

			log.Println("stop signal received! Waiting for active recordings to finish.")
			b.stopActiveProcesses(name)
			wg.Wait()
			log.Println("Stopped!")
			if is_single_run && name == streamerName || b.ListRecorders() == nil {
				b.status.IsRunning = false
				return
			}
		case <-ticker.C:
			// Optionally reload config.
			if config.Settings.AutoReload {
				config.Reload(file.Config_json_path, &config.Settings)
			}

			// Remove finished processes.
			b.checkProcesses()
			// For each streamer in the config, start a recorder if one isn’t already running.
			for _, streamer := range config.Streamers.StreamerList {
				if is_single_run && streamer.Name != streamerName {
					continue
				}
				// Start a new recorder if one isn’t already running.
				wg.Add(1)
				go func(wg *sync.WaitGroup) {
					defer wg.Done()
					if b.isRecorderActive(streamer.Name) {
						return
					}

					// Check if a shutdown is in progress before starting a new recorder.
					select {
					case <-b.ctx.Done():
						break
					default:
					}

					b.runRecordLoop(streamer.Name)

				}(&wg)

			}
			time.Sleep(time.Duration(config.Settings.App.RateLimit.Time) * time.Second)
			if b.isFirstRun {
				b.isFirstRun = false
				ticker.Reset(time.Duration(config.Settings.App.Loop_interval) * time.Minute)

			}
		}
	}
}

// checkProcesses looks through the list of processes and removes any that have finished.
func (b *bot) checkProcesses() {
	b.mux.Lock()
	defer b.mux.Unlock()
	for i := 0; i < len(b.status.Processes); i++ {
		rec := b.status.Processes[i]
		// Use signal 0 to check if process is still running.
		if rec.Cmd == nil || rec.Cmd.Process == nil {
			continue
		}
		if err := rec.Cmd.Process.Signal(syscall.Signal(0)); err != nil {
			log.Printf("[bot]: Process for %s has stopped", rec.Name)
			b.status.Processes = append(b.status.Processes[:i], b.status.Processes[i+1:]...)
			i--
		}
	}
}

// isRecorderActive returns true if a recorder for the given streamer is already running.
func (b *bot) isRecorderActive(streamerName string) bool {
	b.mux.Lock()
	defer b.mux.Unlock()
	for _, rec := range b.status.Processes {
		if rec.Name == streamerName && rec.IsRecording {
			return true
		}
	}
	return false
}

// stopActiveProcesses sends a SIGINT to all active recording processes and waits for them to finish.
func (b *bot) stopActiveProcesses(processName string) {
	stopSingleProcess := false
	if processName != "" {
		stopSingleProcess = true
	}
	b.mux.Lock()
	processesCopy := make([]StreamerStatus, len(b.status.Processes))
	copy(processesCopy, b.status.Processes)
	b.mux.Unlock()

	var wg sync.WaitGroup
	for _, rec := range processesCopy {
		if stopSingleProcess && rec.Name == processName {
			wg.Add(1)
			go stopProcess(&wg, rec)
			break
		} else if stopSingleProcess && rec.Name != processName {
			continue
		}
		wg.Add(1)
		go stopProcess(&wg, rec)
	}
	wg.Wait()
	b.status.IsRunning = false
}
func stopProcess(wg *sync.WaitGroup, rec StreamerStatus) {

	defer wg.Done()
	if rec.Cmd != nil && rec.Cmd.Process != nil {
		log.Printf("[bot]: Stopping recording for %s", rec.Name)
		rec.Cmd.Process.Signal(syscall.SIGINT)
		rec.Cmd.Wait()
	}
}

// writeYoutubeDLConfig writes the youtube-dl configuration file.
func (b *bot) writeYoutubeDLConfig() error {
	// Ensure we start with an empty file.
	f, err := os.Create(file.YoutubeDL_configPath)
	if err != nil {
		return err
	}
	defer f.Close()

	folder := config.Settings.App.Videos_folder
	configLine := fmt.Sprintf("-o \"%s", folder) + "/%(id)s/%(title)s.%(ext)s\""
	_, err = f.Write([]byte(configLine))
	return err
}
