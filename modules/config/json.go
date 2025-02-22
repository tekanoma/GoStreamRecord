package config

import (
	"GoRecordurbate/modules/file"
	"encoding/json"
	"log"
	"os"
)

var (
	Settings  settings
	Streamers = StreamersList{StreamerList: []Streamer{}}
	Users     = Logins{Users: []Login{}}
)

func init() {
	loadConfig(file.Config_json_path, &Settings)
	loadConfig(file.Streamers_json_path, &Streamers)
	loadConfig(file.Users_json_path, &Users)
}
func loadConfig(path string, target any) {
	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			log.Fatalf("Config file %s not found.", path)
		}
		log.Fatalf("Error reading file %s: %v", path, err)
	}

	if ok, err := file.CheckJson(path, target); !ok {
		log.Fatalf("Invalid JSON format in %s: %v", path, err)
	}

	if err = json.Unmarshal(data, target); err != nil {
		log.Fatalf("Failed to parse JSON in %s: %v", path, err)
	}
}

func GenerateDefaultConfig(path string, jsonFile any) {
	f, err := os.Create(path)
	if err != nil {
		log.Fatalf("Failed to create config file: %v", err)
	}
	defer f.Close()
	data, _ := json.MarshalIndent(&jsonFile, "", "  ")
	f.Write(data)
}

func Reload(path string, target any) {
	if err := file.ReadJson(path, target); err != nil {
		log.Printf("Error reloading config: %v", err)
	}
}

func Update(path string, newConfig any) {
	var backup any
	if !readConfig(path, &backup) {
		return
	}
	if !writeConfig(path, newConfig) || !verifyConfig(path, newConfig) {
		writeConfig(path, backup)
	}
}

func verifyConfig(path string, config any) bool {
	ok, err := file.CheckJson(path, config)
	if !ok {
		log.Printf("Config verification failed: %v", err)
		return false
	}
	return true
}

func readConfig(path string, target any) bool {
	if err := file.ReadJson(path, target); err != nil {
		log.Printf("Error reading config: %v", err)
		return false
	}
	return true
}

func writeConfig(path string, config any) bool {
	if err := file.WriteJson(path, config); err != nil {
		log.Printf("Error writing config: %v", err)
		return false
	}
	return true
}
