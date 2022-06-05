package internal

import (
	"image"
	"image_storage/src/internal/domain"
)

type ImageStorage interface {
	Save(myImage *domain.MyImage, data *image.Image) (err error)
	Get(myImage *domain.MyImage) (image *image.Image, err error)
}
