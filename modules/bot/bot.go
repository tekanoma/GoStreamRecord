package bot

import (
	"GoRecordurbate/modules/config"
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

	// processes holds our active recording processes.
	processes []config.Streamer_status

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
	config.C.Reload()

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
	log.Println("Stopping bot..")
	if b.running {
		log.Println("Caught stop signal, stopping")
		b.running = false
	}
	b.processes = []config.Streamer_status{}
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
		log.Printf("Error creating request: %v", err)
		return false
	}
	req.Header.Add("X-Requested-With", "XMLHttpRequest")
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Printf("An error occurred while making the request: %v", err)
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

// IsOnline checks if the streamer is online by sending a GET request to the API.
// It waits for 3 seconds before making the call.
func (b *Bot) IsOnline(username string) bool {
	time.Sleep(3 * time.Second)
	urlStr := "https://chaturbate.com/api/public/affiliates/onlinerooms/?wm=DkfRj&client_ip=request_ip&limit=500"
	urlStr = "https://chaturbate.com/api/chatvideocontext/" + username
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
		log.Printf("Error decoding JSON: %v", err)
		return false
	}

	log.Printf("Checking online for username: %s", username)

	if res.Username == username && res.CurrentShow == "public" {
		return true
	}

	return false
}

// Run starts the main loop of the Bot. While running, it reloads configuration (if enabled),
// checks for finished recording processes, and starts a new recording for any streamer that is online.
func (b *Bot) Run() {

	for b.running {

		// Write youtube dl config
		/*

			Default:
			-o "./videos/%(id)s/%(title)s.%(ext)s"

			To reduce output video filesize, use the following instead to limit to [height<1080][fps<?60]:
			-f 'best[height<1080][fps<?60]' -o "./videos/%(id)s/%(title)s.%(ext)s"
			 --quiet

		*/
		//os.Remove(config.C.YoutubeDL.Config) // ensure empty file
		f, err := os.Create(config.C.YoutubeDL.Config)
		if err != nil {
			log.Println("Error Creating to file: ", err)
			continue
		}
		folder := config.C.App.Videos_folder
		_, err = f.Write([]byte("-o \"" + folder + "/%(id)s/%(title)s.%(ext)s\""))
		if err != nil {
			log.Println("Error writing to file: ", err)
			continue
		}
		// Catch any panic in this iteration.
		func() {
			defer func() {
				if r := recover(); r != nil {
					log.Printf("Loop error: %v", r)
					time.Sleep(1 * time.Second)
				}
			}()

			// Reload config if auto_reload_config is enabled.
			if config.C.AutoReload {
				config.C.Reload()

			}

			// Check current recording processes.
			for i := 0; i < len(b.processes); i++ {
				rec := b.processes[i]
				// Send signal 0 to check if the process is still running.
				err := rec.Cmd.Process.Signal(syscall.Signal(0))
				if err != nil {
					log.Printf("Stopped recording %s", rec.Name)
					// Mark the streamer as not recording.
					for j, s := range config.Streamer_Statuses {
						if s.Name == rec.Name {
							config.Streamer_Statuses[j].Running = false
						}
					}
					// Remove the finished process from our slice.
					b.processes = append(b.processes[:i], b.processes[i+1:]...)
					i-- // adjust index after removal
				}
			}

			var wg sync.WaitGroup
			// Check each streamer and start recording if needed.
			for _, s := range config.C.App.Streamers {
				is_processed := false
				for _, p := range b.processes {
					if p.Name == s.Name {
						is_processed = true
					}
				}

				if is_processed {
					continue
				}

				wg.Add(1)
				go func() {
					defer wg.Done()
					if b.IsOnline(s.Name) {
						log.Printf("Started to record %s", s.Name)
						args := strings.Fields(config.C.YoutubeDL.Binary)
						recordURL := fmt.Sprintf("https://chaturbate.com/%s/", s.Name)
						args = append(args, recordURL, "--config-location", config.C.YoutubeDL.Config)
						cmd := exec.Command(args[0], args[1:]...)

						// Start the recording process.
						if err := cmd.Start(); err != nil {
							log.Printf("Error starting recording for %s: %v", s.Name, err)
						} else {
							b.processes = append(b.processes, config.Streamer_status{Name: s.Name, Cmd: cmd, Running: true})

							for _, p := range b.processes {
								if p.Name == s.Name {
									is_processed = true
								}
							}
							fmt.Println("Recording has started")
						}
					} else {
						log.Printf("Streamer is not online.. %s", s.Name)
					}
					// Respect rate limiting if enabled.
					if config.C.App.RateLimit.Enable {
						time.Sleep(time.Duration(config.C.App.RateLimit.Time) * time.Second)
					}
				}()
				config.C.Update()
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
		for i := range config.Streamer_Statuses {
			config.Streamer_Statuses[i].Running = false
		}
		fmt.Println(rec.Name)
		rec.Cmd.Process.Signal(syscall.SIGINT)
		rec.Cmd.Wait()
	}
	config.C.Update()
	log.Println("Successfully stopped")
}
