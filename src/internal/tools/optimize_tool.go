package tools

import (
	"bytes"
	"image/png"
	"image_storage/src/internal/domain"
	"io/ioutil"
	"strconv"

	"github.com/disintegration/imaging"
)

type OptimizeTool struct{}

func NewOptimizeTool() *OptimizeTool {
	return &OptimizeTool{}
}

func (o *OptimizeTool) ChangeQuality(input *domain.MyImage, quality string) (output *domain.MyImage, err error) {

	img, err := png.Decode(bytes.NewReader(input.Data))
	if err != nil {
		return
	}
	w := new(bytes.Buffer)
	q, err := strconv.Atoi(quality)
	if err != nil {
		return
	}
	err = imaging.Encode(w, img, imaging.JPEG, imaging.JPEGQuality(q))
	if err != nil {
		return
	}
	byteImage, err := ioutil.ReadAll(w)
	output = &domain.MyImage{
		Id:      input.Id,
		Quality: quality,
		Data:    byteImage,
	}
	return
}
