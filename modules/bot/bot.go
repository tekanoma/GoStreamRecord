package bot

import (
	"GoRecordurbate/modules/config"
	"GoRecordurbate/modules/file"
	"context"
	"encoding/json"
	"fmt"
	"io"
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
	isFirstRun bool
	isRunning  bool
	processes  []StreamerStatus
	logger     *log.Logger
	Interface

	// ctx is used to signal shutdown.
	ctx    context.Context
	cancel context.CancelFunc
}

// StreamerStatus holds info for a recording process.
type StreamerStatus struct {
	Name        string
	IsRecording bool
	Cmd         *exec.Cmd
}

// NewBot creates a new Bot, sets up its cancellation context, and registers a signal handler.
func NewBot(logger *log.Logger) *bot {
	ctx, cancel := context.WithCancel(context.Background())
	b := &bot{
		logger:     logger,
		ctx:        ctx,
		cancel:     cancel,
		isRunning:  false,
		isFirstRun: true,
	}
	// Register to catch SIGINT and SIGTERM and trigger Stop.
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		s := <-sigs
		logger.Printf("Caught signal %v, stopping", s)
		b.Stop()
	}()
	return b
}

// AppendStreamer adds a new StreamerStatus for a given streamer.
func (b *bot) ResetBot() *bot {
	b.mux.Lock()
	defer b.mux.Unlock()
	b.Stop()
	b.cancel()
	return NewBot(b.logger)
}

// AppendStreamer adds a new StreamerStatus for a given streamer.
func (b *bot) AppendStreamer(name string) {
	b.mux.Lock()
	defer b.mux.Unlock()
	b.processes = append(b.processes, StreamerStatus{Name: name, IsRecording: false})
}

// ListRecorders returns the current list of recorder statuses.
func (b *bot) ListRecorders() []StreamerStatus {
	b.mux.Lock()
	defer b.mux.Unlock()
	return b.processes
}

// Stop signals the bot to stop starting new recordings and then gracefully stops active processes.
func (b *bot) Stop() {
	// Signal cancellation.
	b.cancel()
	log.Println("Stopping bot..")
	// Give current recorders time to finish (or exit gracefully).
	b.stopActiveProcesses()
}

// IsRoomPublic checks if a given room is public by sending a POST request.
func (b *bot) IsRoomPublic(username string) bool {
	// Wait for the configured rate limit.
	time.Sleep(time.Duration(config.C.App.RateLimit.Time) * time.Second)
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

// IsOnline checks if the streamer is online by sending a GET request.
func (b *bot) IsOnline(username string) bool {
	// Short delay before making the call.
	time.Sleep(3 * time.Second)
	urlStr := "https://chaturbate.com/api/chatvideocontext/" + username
	resp, err := http.Get(urlStr)
	if err != nil {
		log.Printf("Error in GET request: %v", err)
		return false
	}
	defer resp.Body.Close()

	var res struct {
		Username    string `json:"broadcaster_username"`
		CurrentShow string `json:"room_status"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&res); err != nil {
		body, _ := io.ReadAll(resp.Body)
		if err := json.Unmarshal(body, &res); err != nil {
			log.Printf("Error decoding JSON: %v", err)
			log.Println("Most likely rate limit issue..")
			return false
		}
	}
	return res.Username == username && res.CurrentShow == "public"
}

// Run starts the main loop of the Bot.
// It reloads the configuration, checks running processes, and for each configured streamer
// starts a recording if one isn’t already running.
// Once the context is cancelled (via Stop), no new recordings are started.
func (b *bot) Run() {
	b.mux.Lock()
	if b.isRunning {
		b.mux.Unlock()
		return
	}
	b.isRunning = true
	b.mux.Unlock()

	// Write youtube-dl config.
	if err := b.writeYoutubeDLConfig(); err != nil {
		log.Println("Error writing youtube-dl config:", err)
		return
	}

	var wg sync.WaitGroup
	var ticker *time.Ticker
	ticker = time.NewTicker(time.Duration(1) * time.Second)

	defer ticker.Stop()

	// Main loop.
	for {
		select {
		case <-b.ctx.Done():
			log.Println("Shutdown signal received: waiting for active recordings to finish.")
			// Wait for any active record loops to finish.
			wg.Wait()
			// Now send SIGINT to active processes.
			b.stopActiveProcesses()
			return
		case <-ticker.C:
			// Optionally reload config.
			if config.C.AutoReload {
				config.C.Reload()
			}

			// Remove finished processes.
			b.checkProcesses()
			// For each streamer in the config, start a recorder if one isn’t already running.
			for _, streamer := range config.C.App.Streamers {

				if b.isRecorderActive(streamer.Name) {
					continue
				}

				// Check if a shutdown is in progress before starting a new recorder.
				select {
				case <-b.ctx.Done():
					break
				default:
				}

				wg.Add(1)
				go b.runRecordLoop(&wg, streamer.Name)
				// Respect rate limiting.
				time.Sleep(time.Duration(config.C.App.RateLimit.Time) * time.Second)

			}
			if b.isFirstRun {
				b.isFirstRun = false
				fmt.Println(time.Duration(config.C.App.Loop_interval) * time.Minute)
				ticker.Reset(time.Duration(config.C.App.Loop_interval) * time.Minute)

			}
		}
	}
}

// runRecordLoop starts a recording for the given streamer (if online) and waits for the process to finish.
func (b *bot) runRecordLoop(wg *sync.WaitGroup, streamerName string) {
	defer wg.Done()

	// If the bot is stopping, do not check online status.
	select {
	case <-b.ctx.Done():
		return
	default:
	}

	log.Printf("[bot]: Checking %s room status...", streamerName)
	if !b.IsOnline(streamerName) {
		log.Printf("[bot]: Streamer %s is not online.", streamerName)
		return
	}

	log.Printf("[bot]: Starting recording for %s", streamerName)
	args := strings.Fields(config.C.YoutubeDL.Binary)
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
	b.processes = append(b.processes, StreamerStatus{
		Name:        streamerName,
		IsRecording: true,
		Cmd:         cmd,
	})
	b.mux.Unlock()

	// Wait for the command to finish.
	err := cmd.Wait()
	if err != nil {
		log.Printf("[bot]: Recording for %s ended with error: %v", streamerName, err)
		log.Printf("[bot]: Command was: %v", cmd)
	} else {
		log.Printf("[bot]: Recording for %s finished successfully", streamerName)
	}

	// Remove this process from our list.
	b.mux.Lock()
	for i, p := range b.processes {
		if p.Name == streamerName && p.Cmd == cmd {
			b.processes = append(b.processes[:i], b.processes[i+1:]...)
			break
		}
	}
	b.mux.Unlock()
}

// checkProcesses looks through the list of processes and removes any that have finished.
func (b *bot) checkProcesses() {
	b.mux.Lock()
	defer b.mux.Unlock()
	for i := 0; i < len(b.processes); i++ {
		rec := b.processes[i]
		// Use signal 0 to check if process is still running.
		if rec.Cmd == nil || rec.Cmd.Process == nil {
			continue
		}
		if err := rec.Cmd.Process.Signal(syscall.Signal(0)); err != nil {
			log.Printf("[bot]: Process for %s has stopped", rec.Name)
			b.processes = append(b.processes[:i], b.processes[i+1:]...)
			i--
		}
	}
}

// isRecorderActive returns true if a recorder for the given streamer is already running.
func (b *bot) isRecorderActive(streamerName string) bool {
	b.mux.Lock()
	defer b.mux.Unlock()
	for _, rec := range b.processes {
		if rec.Name == streamerName && rec.IsRecording {
			return true
		}
	}
	return false
}

// stopActiveProcesses sends a SIGINT to all active recording processes and waits for them to finish.
func (b *bot) stopActiveProcesses() {
	b.mux.Lock()
	processesCopy := make([]StreamerStatus, len(b.processes))
	copy(processesCopy, b.processes)
	b.mux.Unlock()

	for _, rec := range processesCopy {
		if rec.Cmd != nil && rec.Cmd.Process != nil {
			log.Printf("[bot]: Stopping recording for %s", rec.Name)
			rec.Cmd.Process.Signal(syscall.SIGINT)
			// Wait for process to exit.
			rec.Cmd.Wait()
		}
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

	folder := config.C.App.Videos_folder
	configLine := fmt.Sprintf("-o \"%s", folder) + "/%(id)s/%(title)s.%(ext)s\""
	_, err = f.Write([]byte(configLine))
	return err
}
