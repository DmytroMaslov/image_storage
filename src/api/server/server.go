package server

import (
	"image_storage/src/api/server/controllers"
	"image_storage/src/config"

	"github.com/labstack/echo"
)

type Server struct {
	echo            *echo.Echo
	config          *config.Config
	imageController *controllers.ImageController
}

func NewServer(config *config.Config, imageController *controllers.ImageController) *Server {
	return &Server{
		config:          config,
		imageController: imageController,
	}
}

func (s *Server) Run() (err error) {
	s.echo = echo.New()
	s.initEndpoints()
	err = s.echo.Start(s.config.ServerPort)
	return
}

func (s *Server) initEndpoints() {
	s.echo.POST("/upload/", s.imageController.Upload)
	s.echo.GET("download/:id", s.imageController.Download)
}
