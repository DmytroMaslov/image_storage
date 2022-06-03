package main

import (
	"image_storage/src/api"
	"image_storage/src/config"
)

func main() {
	config := config.GetConfig()
	app := api.NewApp(config)
	app.Run()
}
