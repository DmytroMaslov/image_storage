package localstorage

import (
	"bytes"
	"fmt"
	"image"
	"image/png"
	"image_storage/src/config"
	"image_storage/src/internal"
	"image_storage/src/internal/domain"
	"io/ioutil"
	"log"
	"os"
)

type LocalStorage struct {
	config *config.Config
}

func NewLocalStorage(config *config.Config) internal.ImageStorage {
	var LocalStorage internal.ImageStorage = &LocalStorage{
		config: config,
	}
	return LocalStorage
}

func (ls *LocalStorage) GeneratePath(myImage *domain.MyImage) string {
	//path := "./" + filepath.Join(ls.config.StorageFolderName, myImage.Id, fmt.Sprint(myImage.Quality))
	path := fmt.Sprintf("./%s/%s/%s", ls.config.StorageFolderName, myImage.Id, myImage.Quality)
	return path
}

func (ls *LocalStorage) GeneratePathWithFileName(i *domain.MyImage) string {
	fileName := fmt.Sprintf("/%s_%s.png", i.Id, i.Quality)
	path := ls.GeneratePath(i) + fileName
	return path
}

func (ls *LocalStorage) Save(myImage *domain.MyImage) (err error) {

	img, _, err := image.Decode(bytes.NewReader(myImage.Data))
	if err != nil {
		return
	}
	log.Println(ls.GeneratePath(myImage))
	_, err = os.Stat(ls.GeneratePath(myImage))
	if err != nil {
		log.Println(err)
		err = os.MkdirAll(ls.GeneratePath(myImage), 0755)
		if err != nil {
			log.Println(err)
			return
		}
	}
	out, err := os.Create(ls.GeneratePathWithFileName(myImage))
	if err != nil {
		log.Println(err)
		return
	}
	defer out.Close()
	err = png.Encode(out, img)
	if err != nil {
		return
	}
	return
}

func (ls *LocalStorage) Get(image *domain.MyImage) (err error) {
	path := ls.GeneratePathWithFileName(image)
	image.Data, err = ioutil.ReadFile(path)
	if err != nil {
		return
	}
	return
}
