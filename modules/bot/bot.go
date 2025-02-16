package bot

import (
	"GoRecordurbate/modules/config"
	"GoRecordurbate/modules/file"
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

var (
	START   = "start"
	STOP    = "stop"
	RESTART = "restart"
)

// StreamerStatus tracks a streamer and whether a recording process is running.
type StreamerStatus struct {
	mu        sync.Mutex
	Name      string
	Recording bool
}

// ProcessRecord holds the recording process associated with a streamer.
type ProcessRecord struct {
	streamer string
	cmd      *exec.Cmd
}

// Bot encapsulates the recording bot’s state.
type Bot struct {
	// error flag (unused in this example; set false if no error)
	error bool
	// running flag: set to false when a stop signal is caught
	running bool

	// config holds the configuration loaded from your config module.
	// (Assumes your config module defines a Config struct with fields like AutoReload,
	// YoutubeDL (with Binary and Config), and App with RateLimit and Streamers.)
	config *config.Config

	// streamers holds our local list of streamers along with a “recording” flag.
	streamers []StreamerStatus
	// processes holds our active recording processes.
	processes []ProcessRecord

	logger *log.Logger
}

// NewBot creates a new Bot, loads the configuration, and registers signal handlers.
func NewBot(logger *log.Logger) *Bot {
	b := &Bot{
		error:   false,
		running: true,
		logger:  logger,
	}

	// Load config for the first time.
	b.ReloadConfig()

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

// Stop sets the running flag to false so that the bot can exit gracefully.
func (b *Bot) Stop() {
	if b.running {
		b.logger.Println("Caught stop signal, stopping")
		b.running = false
	}
}

// ReloadConfig loads (or reloads) the configuration.
// On the first call it creates a local list of streamers (all marked as not recording).
// On subsequent calls it removes streamers no longer in the config and adds any new ones.
func (b *Bot) ReloadConfig() {
	var newConfig config.Config
	err := file.ReadJson("./config.json", &newConfig)
	if err != nil {
		b.logger.Printf("Error loading config: %v", err)
		return
	}

	// If this is the first load, initialize our streamer list.
	if b.config == nil {
		b.config = &newConfig
		for _, s := range b.config.App.Streamers {
			b.streamers = append(b.streamers, StreamerStatus{Name: s.Name, Recording: false})
		}
		return
	}

	// Remove streamers that were removed in the new config.
	newStreamerMap := make(map[string]bool)
	for _, s := range newConfig.App.Streamers {
		newStreamerMap[s.Name] = true
	}
	updatedStreamers := []StreamerStatus{}
	for _, s := range b.streamers {
		if _, exists := newStreamerMap[s.Name]; exists {
			updatedStreamers = append(updatedStreamers, s)
		} else {
			b.logger.Printf("%s has been removed", s.Name)
		}
	}
	b.streamers = updatedStreamers

	// Add any new streamers.
	currentMap := make(map[string]bool)
	for _, s := range b.streamers {
		currentMap[s.Name] = true
	}
	for _, s := range newConfig.App.Streamers {
		if !currentMap[s.Name] {
			b.streamers = append(b.streamers, StreamerStatus{Name: s.Name, Recording: false})
		}
	}
	b.config = &newConfig
}

// IsRoomPublic checks if a given room (streamer) is public by sending a POST request.
// It waits for 3 seconds before making the call.
func (b *Bot) IsRoomPublic(username string) bool {
	time.Sleep(3 * time.Second)
	urlStr := "https://chaturbate.com/get_edge_hls_url_ajax/"
	data := url.Values{}
	data.Set("room_slug", username)
	req, err := http.NewRequest("POST", urlStr, strings.NewReader(data.Encode()))
	if err != nil {
		b.logger.Printf("Error creating request: %v", err)
		return false
	}
	req.Header.Add("X-Requested-With", "XMLHttpRequest")
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		b.logger.Printf("An error occurred while making the request: %v", err)
		return false
	}
	defer resp.Body.Close()

	var res struct {
		Success    bool   `json:"success"`
		RoomStatus string `json:"room_status"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&res); err != nil {
		b.logger.Printf("Error decoding JSON: %v", err)
		return false
	}
	return res.Success && res.RoomStatus == "public"
}

// IsOnline checks if the streamer is online by sending a GET request to the API.
// It waits for 3 seconds before making the call.
func (b *Bot) IsOnline(username string) bool {
	time.Sleep(3 * time.Second)
	urlStr := "https://chaturbate.com/api/public/affiliates/onlinerooms/?wm=DkfRj&client_ip=request_ip&limit=500"
	urlStr = "https://chaturbate.com/api/chatvideocontext/" + username
	resp, err := http.Get(urlStr)
	if err != nil {
		b.logger.Printf("Error in GET request: %v", err)
		return false
	}
	defer resp.Body.Close()

	var res struct {
		Username    string `json:"broadcaster_username"`
		CurrentShow string `json:"room_status"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&res); err != nil {
		b.logger.Printf("Error decoding JSON: %v", err)
		return false
	}

	b.logger.Printf("Checking online for username: %s", username)

	if res.Username == username && res.CurrentShow == "public" {
		return true
	}

	return false
}

// Run starts the main loop of the Bot. While running, it reloads configuration (if enabled),
// checks for finished recording processes, and starts a new recording for any streamer that is online.
func (b *Bot) Run() {
	for b.running {
		// Catch any panic in this iteration.
		func() {
			defer func() {
				if r := recover(); r != nil {
					b.logger.Printf("Loop error: %v", r)
					time.Sleep(1 * time.Second)
				}
			}()

			// Reload config if auto_reload_config is enabled.
			if b.config != nil && b.config.AutoReload {
				b.ReloadConfig()
			}

			// Check current recording processes.
			for i := 0; i < len(b.processes); i++ {
				rec := b.processes[i]
				// Send signal 0 to check if the process is still running.
				err := rec.cmd.Process.Signal(syscall.Signal(0))
				if err != nil {
					b.logger.Printf("Stopped recording %s", rec.streamer)
					// Mark the streamer as not recording.
					for j, s := range b.streamers {
						if s.Name == rec.streamer {
							b.streamers[j].Recording = false
						}
					}
					// Remove the finished process from our slice.
					b.processes = append(b.processes[:i], b.processes[i+1:]...)
					i-- // adjust index after removal
				}
			}
			var wg sync.WaitGroup
			// Check each streamer and start recording if needed.
			for i, s := range b.streamers {
				s.mu.Lock()
				defer s.mu.Unlock()
				if s.Recording {
					b.logger.Println("Recorder already running!")
					continue
				}
				wg.Add(1)
				go func() {
					defer wg.Done()
					if b.IsOnline(s.Name) {
						b.logger.Printf("Started to record %s", s.Name)
						// Prepare the command.
						// We assume your configuration holds the YouTube-DL command as a string in YoutubeDL.Binary
						// and the config location in YoutubeDL.Config.
						args := strings.Fields(b.config.YoutubeDL.Binary)
						recordURL := fmt.Sprintf("https://chaturbate.com/%s/", s.Name)
						args = append(args, recordURL, "--config-location", b.config.YoutubeDL.Config)
						cmd := exec.Command(args[0], args[1:]...)
						log.Println(cmd)
						// Start the recording process.
						if err := cmd.Start(); err != nil {
							b.logger.Printf("Error starting recording for %s: %v", s.Name, err)
						} else {
							b.processes = append(b.processes, ProcessRecord{streamer: s.Name, cmd: cmd})
							b.streamers[i].Recording = true
							fmt.Println("Recording has started")
						}
					} else {
						b.logger.Printf("Streamer is not online.. %s", s.Name)
					}
					// Respect rate limiting if enabled.
					if b.config.App.RateLimit.Enable {
						time.Sleep(time.Duration(b.config.App.RateLimit.Time) * time.Second)
					}
				}()
				time.Sleep(time.Duration(config.C.App.RateLimit.Time) * time.Second)
			}
			// Wait for 1 minute in 1-second intervals.
			for i := 0; i < 60; i++ {
				if !b.running {
					break
				}
				time.Sleep(1 * time.Second)
			}
		}()
	}

	// When the loop ends, stop all active recording processes.
	for _, rec := range b.processes {
		rec.cmd.Process.Signal(syscall.SIGINT)
		rec.cmd.Wait()
	}
	b.logger.Println("Successfully stopped")
}
