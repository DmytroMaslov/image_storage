package mock

import (
	"image"
	"image_storage/src/internal/domain"
)

type ImageStorageMock struct {
	FuncSave func(myImage *domain.MyImage, data *image.Image) (err error)
	FuncGet  func(myImage *domain.MyImage) (image *image.Image, err error)
}

func (is *ImageStorageMock) Save(myImage *domain.MyImage, data *image.Image) (err error) {
	return is.FuncSave(myImage, data)
}

func (is *ImageStorageMock) Get(myImage *domain.MyImage) (image *image.Image, err error) {
	return is.FuncGet(myImage)
}
