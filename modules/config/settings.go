package config

import (
	"sync"
)

type settings struct {
	App        app       `json:"app"`
	YoutubeDL  youtubeDL `json:"youtube-dl"`
	AutoReload bool      `json:"auto_reload_config"`
}

type youtubeDL struct {
	Binary string `json:"binary"`
}
type app struct {
	mux           sync.Mutex
	Port          int        `json:"port"`
	Loop_interval int        `json:"loop_interval_in_minutes"`
	Videos_folder string     `json:"video_output_folder"`
	RateLimit     rate_limit `json:"rate_limit"`
	ExportPath    string     `json:"default_export_location"`
}
type rate_limit struct {
	Enable bool `json:"enable"`
	Time   int  `json:"time"`
}
