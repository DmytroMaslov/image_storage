package controllers

import (
	"bytes"
	"fmt"
	"image"
	"image_storage/src/internal/services"
	"image_storage/src/pkg"
	"net/http"
	"strconv"

	"github.com/labstack/echo"
)

type ImageController struct {
	imageService services.ImageUseCases
	log          pkg.Logger
}

func NewImageController(imageService services.ImageUseCases, logger pkg.Logger) *ImageController {
	return &ImageController{
		imageService: imageService,
		log:          logger,
	}
}

func (c *ImageController) Upload(e echo.Context) error {
	file, err := e.FormFile("file")
	if err != nil {
		c.log.Errorf("can't get file from request, reason: %s", err.Error())
		return e.String(http.StatusBadRequest, "provide file")
	}
	src, err := file.Open()
	if err != nil {
		c.log.Errorf("can't open file from request, reason: %s", err.Error())
		return e.String(http.StatusBadRequest, "provide valid file format")
	}
	defer src.Close()
	img, _, err := image.Decode(src)
	if err != nil {
		c.log.Errorf("can't decode file from request, reason: %s", err.Error())
		return e.String(http.StatusBadRequest, "provide valid file format")
	}
	myImg, err := c.imageService.UploadImage(&img)
	if err != nil {
		c.log.Errorf("can't save file from request, reason: %s", err.Error())
		return e.String(http.StatusInternalServerError, "file not saved")
	}
	return e.String(http.StatusOK, fmt.Sprintf("image id: %s", myImg.Id))
}

func (c *ImageController) Download(e echo.Context) error {
	id := e.Param("id")
	if id == "" {
		return e.String(http.StatusBadRequest, "provide valid image id")
	}
	qualityString := e.QueryParam("quality")
	if qualityString == "" {
		return e.String(http.StatusBadRequest, "provide valid quality")
	}
	quality, err := strconv.Atoi(qualityString)
	if err != nil {
		return e.String(http.StatusBadRequest, "provide valid quality")
	}
	byteImage, err := c.imageService.GetByteImageByIdAndQuality(id, quality)
	if err != nil {
		return e.String(http.StatusNotFound, "can't fiend image by given params")
	}
	return e.Stream(http.StatusOK, "image/jpeg", bytes.NewReader(*byteImage))

}
