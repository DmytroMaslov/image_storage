package main

import (
	"image_storage/src/api"
	"image_storage/src/config"
	"image_storage/src/pkg"
)

func main() {
	config := config.GetConfig()
	logger := pkg.NewLogger()
	app := api.NewApp(config, logger)
	app.Run()
}
