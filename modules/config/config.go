package config

import (
	"GoRecordurbate/modules/file"
	"fmt"
	"log"
	"os"
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
}
type app struct {
	mux sync.Mutex

	Port          int        `json:"port"`
	Videos_folder string     `json:"video_output_folder"`
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

var (
	C Config
)

func (c *Config) Init() {
	_, err := os.ReadFile(file.Config_path)
	if os.IsNotExist(err) {
		log.Fatal(file.Config_path, " was not found..")
	}

	ok, err := file.CheckJson(file.Config_path, &c)
	if !ok {
		log.Fatal("Config error: ", err)
		return
	}
	err = file.ReadJson(file.Config_path, &c)
	if err != nil {
		log.Fatal("Error reading config: ", err)
		return
	}
}

func (c *Config) Reload() {

	err := file.ReadJson(file.Config_path, &c)
	if err != nil {
		log.Printf("Error loading config: %v", err)
		return
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
	ok, err := file.CheckJson(file.Config_path, &c)
	if !ok {
		log.Fatal("Config error: ", err)
		return false
	}
	return true
}

func (c *Config) read() bool {
	err := file.ReadJson(file.Config_path, &c)
	if err != nil {
		log.Println("Error reading config: ", err)
		return false
	}
	return true
}

func (c *Config) write() bool {
	err := file.WriteJson(file.Config_path, &c)
	if err != nil {
		log.Println("Error writing config: ", err)
		return false
	}
	return true
}

func (app *app) AddStreamer(streamerName string) string {
	app.mux.Lock()
	defer app.mux.Unlock()
	for _, streamer := range app.Streamers {
		if streamerName == streamer.Name {
			log.Printf("%s has already been addded.", streamerName)
			return fmt.Sprintf("%s has already been addded.", streamerName)
		}
	}
	app.Streamers = append(app.Streamers, streamer{Name: streamerName})
	C.App = *app
	ok := C.write()
	if !ok {
		log.Printf("Error adding %s..\n", streamerName)
		return fmt.Sprintf("Error adding %s..\n", streamerName)
	}
	log.Printf("%s has been added", streamerName)
	return fmt.Sprintf("%s has been added", streamerName)
}

func (app *app) RemoveStreamer(streamerName string) string {
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
		return fmt.Sprintf("%s hasn't been added", streamerName)
	}
	app.Streamers = newList
	C.App = *app
	ok := C.write()
	if !ok {
		log.Printf("Error removing %s..\n", streamerName)
		return fmt.Sprintf("Error removing %s..\n", streamerName)
	}
	log.Printf("%s has been deleted", streamerName)
	return fmt.Sprintf("%s has been deleted", streamerName)

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
