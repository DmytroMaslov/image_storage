package services

import (
	"encoding/json"
	"image"
	"image_storage/src/config"
	"image_storage/src/internal"
	"image_storage/src/internal/domain"
	"image_storage/src/internal/tools"
	"image_storage/src/pkg"
	"log"

	"github.com/google/uuid"
)

type ImageService struct {
	producer internal.ImageProducer
	convert  *tools.ConvertTool
	repo     internal.ImageStorage
	log      pkg.Logger
	config   *config.Config
}
type ImageUseCases interface {
	UploadImage(image *image.Image) (myImage *domain.MyImage, err error)
	GetByteImageByIdAndQuality(id string, quality int) (byteImage *[]byte, err error)
	ReduceQuality(myImage *domain.MyImage, quality []int) (err error)
}

func NewImageService(producer internal.ImageProducer, optimize *tools.ConvertTool, repo internal.ImageStorage, logger pkg.Logger, config *config.Config) ImageUseCases {
	var ImageService ImageUseCases = &ImageService{
		producer: producer,
		convert:  optimize,
		repo:     repo,
		log:      logger,
		config:   config,
	}
	return ImageService
}

func (s *ImageService) UploadImage(image *image.Image) (myImage *domain.MyImage, err error) {
	myImage = &domain.MyImage{
		Id:      uuid.New().String(),
		Quality: s.config.InitialQuality,
	}
	err = s.Save(myImage, image)
	if err != nil {
		return
	}
	s.log.Infof("image successful saved with id: %s quality %v ", myImage.Id, myImage.Quality)
	err = s.SendImageToQueue(myImage)
	if err != nil {
		return
	}
	s.log.Infof("image successful sended to queue with id: %s", myImage.Id)
	return

}

func (s *ImageService) SendImageToQueue(image *domain.MyImage) error {
	byteImage, err := json.Marshal(image)
	if err != nil {
		log.Println(err)
		return err
	}
	err = s.producer.Publish(byteImage, "application/json")
	if err != nil {
		return err
	}
	return nil
}

func (s *ImageService) convertImageToJpeg(input *image.Image) (output *image.Image, err error) {
	convertedImage, err := s.convert.ConvertImageToJpeg(input)
	if err != nil {
		return
	}
	output = &convertedImage
	return
}

func (s *ImageService) Save(myImage *domain.MyImage, image *image.Image) (err error) {
	convertedImage, err := s.convertImageToJpeg(image)
	if err != nil {
		return
	}
	err = s.repo.Save(myImage, convertedImage)
	return
}

func (s *ImageService) GetByteImageByIdAndQuality(id string, quality int) (byteImage *[]byte, err error) {
	newImage := &domain.MyImage{
		Id:      id,
		Quality: quality,
	}
	img, err := s.repo.Get(newImage)
	if err != nil {
		return
	}
	s.log.Infof("successful find image with id:%s quality:%v", id, quality)
	return s.convert.ConvertImageToByteArray(img)
}

func (s *ImageService) ReduceQuality(myImage *domain.MyImage, quality []int) (err error) {
	img, err := s.repo.Get(myImage)
	if err != nil {
		return
	}
	for _, q := range quality {
		myImage.Quality = q
		err = s.Save(myImage, img)
		if err != nil {
			return
		}
		s.log.Infof("successful reduce quality for image id:%s, quality: %v", myImage.Id, myImage.Quality)
	}
	return
}
