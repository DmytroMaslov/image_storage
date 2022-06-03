package api

import (
	"image_storage/src/api/server"
	"image_storage/src/api/server/controllers"
	"image_storage/src/config"
	"image_storage/src/internal/repository/message_broker/rabbitmq"
	localstorage "image_storage/src/internal/repository/storage/local_storage"
	"image_storage/src/internal/services"
	"image_storage/src/internal/tools"
	"log"
)

type App struct {
	config *config.Config
}

func NewApp(config *config.Config) *App {
	return &App{
		config: config,
	}
}

func (a *App) Run() {
	optimizeTool := tools.NewOptimizeTool()
	storage := localstorage.NewLocalStorage(a.config)

	var err error
	producer, err := rabbitmq.NewImageProducer(a.config)
	if err != nil {
		log.Fatal(err)
	}
	imageService := services.NewImageService(producer, optimizeTool, storage)
	consumer, err := rabbitmq.NewImageConsumer(a.config, imageService)
	if err != nil {
		log.Fatal(err)
	}
	err = consumer.RunConsumer(1)
	if err != nil {
		log.Fatal(err)
	}
	controller := controllers.NewImageController(imageService)
	server := server.NewServer(a.config, controller)
	err = server.Run()
	if err != nil {
		log.Fatal(err)
	}
}
