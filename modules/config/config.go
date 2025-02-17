package config

import (
	"GoRecordurbate/modules/file"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"
	"sync"
)

type Config struct {
	App        app       `json:"app"`
	YoutubeDL  youtubeDL `json:"youtube-dl"`
	AutoReload bool      `json:"auto_reload_config"`
}

type youtubeDL struct {
	Binary string `json:"binary"`
	Config string `json:"config"`
}
type app struct {
	mux sync.Mutex

	Port          int        `json:"port"`
	Videos_folder string     `json:"output_folder"`
	RateLimit     rate_limit `json:"rate_limit"`
	ExportPath    string     `json:"default_export_location"`
	Streamers     []streamer `json:"streamers"`
}
type rate_limit struct {
	Enable bool `json:"enable"`
	Time   int  `json:"time"`
}
type streamer struct {
	Name string `json:"name"`
}

// internal for running recorders. not saved on shutdown.
type Streamer_status struct {
	Name    string
	Running bool
	Cmd     *exec.Cmd
}

var (
	Streamer_Statuses = []Streamer_status{}
	C                 Config
	Path              string
)

func (c *Config) Init(config_path string) {
	Path = config_path
	if len(Path) == 0 {
		log.Fatal("Config path not provided")
	}
	_, err := os.ReadFile(Path)
	if os.IsNotExist(err) {
		log.Fatal(Path, " was not found..")
	}

	ok, err := file.CheckJson(Path, &c)
	if !ok {
		log.Fatal("Config error: ", err)
		return
	}
	err = file.ReadJson(Path, &c)
	if err != nil {
		log.Fatal("Error reading config: ", err)
		return
	}
}

func (c *Config) Reload() {
	var newConfig Config
	err := file.ReadJson(Path, &newConfig)
	if err != nil {
		log.Printf("Error loading config: %v", err)
		return
	}

	// If this is the first load, initialize our streamer list.
	if c == nil {
		c = &newConfig
		for _, s := range c.App.Streamers {
			Streamer_Statuses = append(Streamer_Statuses, Streamer_status{Name: s.Name, Running: false})
		}
		return
	}

	// Remove streamers that were removed in the new config.
	newStreamerMap := make(map[string]bool)
	for _, s := range newConfig.App.Streamers {
		newStreamerMap[s.Name] = true
	}
	updatedStreamers := []streamer{}
	for _, s := range c.App.Streamers {
		if _, exists := newStreamerMap[s.Name]; exists {
			updatedStreamers = append(updatedStreamers, s)
		} else {
			log.Printf("%s has been removed", s.Name)
		}
	}
	c.App.Streamers = updatedStreamers
	// Add and log any new streamers.
	currentMap := make(map[string]bool)
	for _, s := range c.App.Streamers {
		currentMap[s.Name] = true
	}
	for _, s := range newConfig.App.Streamers {
		if !currentMap[s.Name] {
			Streamer_Statuses = append(Streamer_Statuses, Streamer_status{Name: s.Name, Running: false})
			log.Printf("%s has been added", s.Name)
		}
	}
}

func (c *Config) Update() {
	var tmpConfig Config
	tmpConfig.read()
	c.write()
	if !c.verify() {
		tmpConfig.write()
	}
}

func (c *Config) verify() bool {
	ok, err := file.CheckJson(Path, &c)
	if !ok {
		log.Fatal("Config error: ", err)
		return false
	}
	return true
}

func (c *Config) read() bool {
	err := file.ReadJson(Path, &c)
	if err != nil {
		log.Println("Error reading config: ", err)
		return false
	}
	return true
}

func (c *Config) write() bool {
	err := file.WriteJson(Path, &c)
	if err != nil {
		log.Println("Error writing config: ", err)
		return false
	}
	return true
}

func (app *app) AddStreamer(streamerName string) {
	app.mux.Lock()
	defer app.mux.Unlock()
	for _, streamer := range app.Streamers {
		if streamerName == streamer.Name {
			log.Printf("%s has already been addded.", streamerName)
			return
		}
	}
	app.Streamers = append(app.Streamers, streamer{Name: streamerName})
	C.App = *app
	ok := C.write()
	if !ok {
		log.Printf("Error adding %s..\n", streamerName)
	}
	log.Printf("%s has been added", streamerName)
}

func (app *app) RemoveStreamer(streamerName string) {
	app.mux.Lock()
	defer app.mux.Unlock()
	newList := []streamer{}
	var wasAdded bool
	for _, streamer := range app.Streamers {
		if streamerName == streamer.Name {
			wasAdded = true
			continue
		}
		newList = append(newList, streamer)
	}
	if !wasAdded {
		log.Printf("%s hasn't been added", streamerName)
		return
	}
	app.Streamers = newList
	C.App = *app
	ok := C.write()
	if !ok {
		log.Printf("Error removing %s..\n", streamerName)
	}
	log.Printf("%s has been deleted", streamerName)

}

func (app *app) ListStreamers() {
	app.mux.Lock()
	defer app.mux.Unlock()
	fmt.Println("Streamers in recording list:")
	for _, streamer := range app.Streamers {
		fmt.Printf("- %s", streamer.Name)
	}
}
func (app *app) ImportStreamers(importFile string) {
	app.mux.Lock()
	defer app.mux.Unlock()
	fileContent, err := os.ReadFile(importFile)
	if os.IsNotExist(err) {
		fmt.Println("File dont exist!")
	}

	for _, line := range strings.Split(string(fileContent), "\n") {
		for _, streamer := range app.Streamers {
			if line == streamer.Name {
				fmt.Printf("%s has already been added", streamer.Name)
				continue
			}
			app.Streamers = append(app.Streamers, streamer)
		}
	}

	C.App = *app
	ok := C.write()
	if !ok {
		log.Printf("Error importing from %s..\n", importFile)
	}
	fmt.Println("Streamers imported!")
}

func (app *app) ExportStreamers() {

	app.mux.Lock()
	defer app.mux.Unlock()
	os.Create(app.ExportPath)
	file, err := os.Open(app.ExportPath)
	if err != nil {
		fmt.Println("Error exporting streamers!")
	}
	defer file.Close()

	for _, streamer := range app.Streamers {
		file.Write([]byte(fmt.Sprintf("%s\n", streamer.Name)))
	}

	fmt.Println("Streamers exported to", app.ExportPath)
}
