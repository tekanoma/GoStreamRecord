package config

import (
	"GoRecordurbate/modules/file"
	"fmt"
	"log"
	"os"
	"strings"
)

type StreamersList struct {
	StreamerList []Streamer
}

type Streamer struct {
	Name string `json:"name"`
}

func AddStreamer(streamerName string) string {
	for _, streamer := range Streamers.StreamerList {
		if streamerName == streamer.Name {
			log.Printf("%s has already been addded.", streamerName)
			return fmt.Sprintf("%s has already been addded.", streamerName)
		}
	}
	Streamers.StreamerList = append(Streamers.StreamerList, Streamer{Name: streamerName})
	ok := writeConfig(file.Streamers_json_path, Streamers)
	if !ok {
		log.Printf("Error adding %s..\n", streamerName)
		return fmt.Sprintf("Error adding %s..\n", streamerName)
	}
	log.Printf("%s has been added", streamerName)
	return fmt.Sprintf("%s has been added", streamerName)
}

func appendStreamer(newList []string) {
	for _, line := range newList {
		exist := false
		for _, streamer := range Streamers.StreamerList {
			if line == streamer.Name {
				fmt.Printf("%s has already been added", streamer)
				exist = true
				break
			}

		}
		if exist {
			continue
		}
		Streamers.StreamerList = append(Streamers.StreamerList, Streamer{Name: line})
	}
}
func RemoveStreamer(streamerName string) string {
	newList := []Streamer{}
	var wasAdded bool
	for _, streamer := range Streamers.StreamerList {
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
	Streamers.StreamerList = newList
	ok := writeConfig(file.Streamers_json_path, Streamers)
	if !ok {
		log.Printf("Error removing %s..\n", streamerName)
		return fmt.Sprintf("Error removing %s..\n", streamerName)
	}
	log.Printf("%s has been deleted", streamerName)
	return fmt.Sprintf("%s has been deleted", streamerName)

}

// TODO: Send to frontend webUI
func (s *Streamer) ListStreamers() {
	fmt.Println("Streamers in recording list:")
	log.Println("Streamers in recording list:")
	for _, streamer := range Streamers.StreamerList {
		fmt.Printf("- %s", streamer)
		log.Printf("- %s", streamer)
	}
}
