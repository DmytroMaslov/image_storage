package services

import (
	"image"
	"image_storage/src/config"
	"image_storage/src/internal/domain"
	"image_storage/src/internal/mock"
	"image_storage/src/pkg"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_UploadImage(t *testing.T) {
	producer := mock.ProducerMock{}
	convert := mock.ConverterToolMock{}
	repo := mock.ImageStorageMock{}
	log := pkg.NewLogger()
	log.InitLogger()
	config := config.GetConfig()
	img, err := GetTestImage()
	assert.Nil(t, err)
	service := NewImageService(&producer, &convert, &repo, log, config)
	repo.FuncSave = func(*domain.MyImage, *image.Image) error {
		return nil
	}
	producer.FunctionPublish = func([]byte, string) error {
		return nil
	}
	convert.FuncConvertImageToJpeg = func(input *image.Image) (output image.Image, err error) {
		return *input, nil
	}
	myImg, err := service.UploadImage(img)
	assert.Nil(t, err)
	assert.NotEmpty(t, myImg)
	assert.Equal(t, config.InitialQuality, myImg.Quality)

}

func GetTestImage() (*image.Image, error) {

	filePath := "./test_data/sample.jpeg"
	f, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	image, _, err := image.Decode(f)
	return &image, err
}
