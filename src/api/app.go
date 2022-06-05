package api

import (
	"image_storage/src/api/server"
	"image_storage/src/api/server/controllers"
	"image_storage/src/config"
	"image_storage/src/internal/repository/message_broker/rabbitmq"
	localstorage "image_storage/src/internal/repository/storage/local_storage"
	"image_storage/src/internal/services"
	"image_storage/src/internal/tools"
	"image_storage/src/pkg"
)

type App struct {
	config *config.Config
	log    pkg.Logger
}

func NewApp(config *config.Config, logger pkg.Logger) *App {
	return &App{
		config: config,
		log:    logger,
	}
}

func (a *App) Run() {
	a.log.InitLogger()
	optimizeTool := tools.NewOptimizeTool(a.log)
	storage := localstorage.NewLocalStorage(a.config, a.log)
	var err error
	producer, err := rabbitmq.NewImageProducer(a.config, a.log)
	if err != nil {
		a.log.Fatalf(err.Error())
	}
	imageService := services.NewImageService(producer, optimizeTool, storage, a.log, a.config)
	consumer, err := rabbitmq.NewImageConsumer(a.config, imageService, a.log)
	if err != nil {
		a.log.Fatalf(err.Error())
	}
	err = consumer.RunConsumer(a.config.WorkerPoolSize, a.config.QualityArray)
	if err != nil {
		a.log.Fatalf(err.Error())
	}
	controller := controllers.NewImageController(imageService, a.log)
	server := server.NewServer(a.config, controller)
	err = server.Run()
	if err != nil {
		a.log.Fatalf(err.Error())
	}
}
