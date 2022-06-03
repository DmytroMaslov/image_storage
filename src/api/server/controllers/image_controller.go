package controllers

import (
	"bytes"
	"image_storage/src/internal/domain"
	"image_storage/src/internal/services"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/google/uuid"
	"github.com/labstack/echo"
)

type ImageController struct {
	imageService services.ImageUseCases
}

func NewImageController(imageService services.ImageUseCases) *ImageController {
	return &ImageController{
		imageService: imageService,
	}
}

func (c *ImageController) Upload(e echo.Context) error {
	file, err := e.FormFile("file")
	if err != nil {
		return e.String(http.StatusInternalServerError, err.Error())
	}
	src, err := file.Open()
	if err != nil {
		return e.String(http.StatusInternalServerError, err.Error())
	}
	defer src.Close()
	byteImageData, err := ioutil.ReadAll(src)
	if err != nil {
		return e.String(http.StatusInternalServerError, err.Error())
	}
	image := &domain.MyImage{
		Id:      uuid.New().String(),
		Quality: "100",
		Data:    byteImageData,
	}
	err = c.imageService.Save(image)
	if err != nil {
		return e.String(http.StatusInternalServerError, err.Error())
	}
	err = c.imageService.SendImageToQueue(image)
	if err != nil {
		return err
	}
	return e.String(http.StatusOK, image.Id)
}

func (c *ImageController) Download(e echo.Context) error {
	id := e.Param("id")
	quality := e.QueryParam("quality")
	log.Printf("id: %v, quality: %v", id, quality)
	myImage, err := c.imageService.GetByIdAndQuality(id, quality)
	if err != nil {
		return err
	}
	return e.Stream(http.StatusOK, "image/png", bytes.NewReader(myImage.Data))

}
