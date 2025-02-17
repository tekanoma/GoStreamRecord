package main

import (
	"GoRecordurbate/modules/config"
	"GoRecordurbate/modules/file"
	"GoRecordurbate/modules/handlers"
	"encoding/json"
	"fmt"
	"os"
)

func init() {

	file.Config_path = "./output/app_config.json"
	_, err := os.Stat(file.Config_path)
	if os.IsNotExist(err) {
		fmt.Println("No config found. Generating.. Please fill in details in:", file.Config_path)
		f, _ := os.Create(file.Config_path)
		tmp := config.Config{}
		b, _ := json.Marshal(tmp)
		f.Write(b)
		os.Exit(0)
	}
	file.InitLog("./output/app.log")
	config.C.Init("./output/app_config.json")

}

// bot/command.go handles incoming commands
func main() {
	handlers.StartWebUI()
}
