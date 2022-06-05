package localstorage

import (
	"fmt"
	"image"
	"image/jpeg"
	"image_storage/src/config"
	"image_storage/src/internal"
	"image_storage/src/internal/domain"
	"image_storage/src/internal/image_errors"
	"image_storage/src/pkg"
	"os"
)

type LocalStorage struct {
	config *config.Config
	log    pkg.Logger
}

func NewLocalStorage(config *config.Config, logger pkg.Logger) internal.ImageStorage {
	var LocalStorage internal.ImageStorage = &LocalStorage{
		config: config,
		log:    logger,
	}
	return LocalStorage
}

func (ls *LocalStorage) GeneratePath(myImage *domain.MyImage) string {
	path := fmt.Sprintf("./%s/%s/%v", ls.config.StorageFolderName, myImage.Id, myImage.Quality)
	return path
}

func (ls *LocalStorage) GeneratePathWithFileName(i *domain.MyImage) string {
	fileName := fmt.Sprintf("/%s_%v.jpeg", i.Id, i.Quality)
	path := ls.GeneratePath(i) + fileName
	return path
}

func (ls *LocalStorage) Save(myImage *domain.MyImage, data *image.Image) (err error) {
	_, err = os.Stat(ls.GeneratePath(myImage))
	if err != nil {
		err = os.MkdirAll(ls.GeneratePath(myImage), 0755)
		if err != nil {
			ls.log.Errorf("local storage can't create folder, reason: %s", err.Error())
			return image_errors.ErrFolderCreate
		}
	}
	out, err := os.Create(ls.GeneratePathWithFileName(myImage))
	if err != nil {
		ls.log.Errorf("local storage can't create file, reason: %s", err.Error())
		return image_errors.ErrFileCreate
	}
	defer out.Close()
	err = jpeg.Encode(out, *data, &jpeg.Options{Quality: myImage.Quality})
	if err != nil {
		ls.log.Errorf("local storage can't save image to file, reason: %s", err.Error())
		return image_errors.ErrSaveToFile
	}
	return
}

func (ls *LocalStorage) Get(myImage *domain.MyImage) (image *image.Image, err error) {
	path := ls.GeneratePathWithFileName(myImage)
	f, err := os.Open(path)
	if err != nil {
		ls.log.Errorf("local storage can't open file, reason: %s", err.Error())
		err = image_errors.ErrFileOpen
		return
	}
	imageFromStorage, err := jpeg.Decode(f)
	if err != nil {
		ls.log.Errorf("local storage can't decode image, reason: %s", err.Error())
		err = image_errors.ErrImageDecode
		return
	}
	image = &imageFromStorage
	return
}
