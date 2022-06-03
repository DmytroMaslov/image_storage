package services

import (
	"encoding/json"
	"image_storage/src/internal"
	"image_storage/src/internal/domain"
	"image_storage/src/internal/tools"
)

type ImageService struct {
	producer internal.ImageProducer
	optimize *tools.OptimizeTool
	repo     internal.ImageStorage
}
type ImageUseCases interface {
	SendImageToQueue(image *domain.MyImage) error
	OptimizeImages(myImage *domain.MyImage) (newImage *domain.MyImage, err error)
	Save(myImage *domain.MyImage) (err error)
	GetByIdAndQuality(id string, quality string) (newImage *domain.MyImage, err error)
}

func NewImageService(producer internal.ImageProducer, optimize *tools.OptimizeTool, repo internal.ImageStorage) ImageUseCases {
	var ImageService ImageUseCases = &ImageService{
		producer: producer,
		optimize: optimize,
		repo:     repo,
	}
	return ImageService
}

func (s *ImageService) SendImageToQueue(image *domain.MyImage) error {
	byteImage, err := json.Marshal(image)
	if err != nil {
		return nil
	}
	err = s.producer.Publish(byteImage, "application/json")
	if err != nil {
		return err
	}
	return nil
}

func (s *ImageService) OptimizeImages(myImage *domain.MyImage) (newImage *domain.MyImage, err error) {
	newImage, err = s.optimize.Optimize(myImage, "80")
	return
}

func (s *ImageService) Save(myImage *domain.MyImage) (err error) {
	err = s.repo.Save(myImage)
	return
}

func (s *ImageService) GetByIdAndQuality(id string, quality string) (newImage *domain.MyImage, err error) {
	newImage = &domain.MyImage{
		Id:      id,
		Quality: quality,
	}
	err = s.repo.Get(newImage)
	return
}
