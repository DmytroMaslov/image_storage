package tools

import (
	"bytes"
	"image"
	"image/jpeg"
	"image_storage/src/internal/image_errors"
	"image_storage/src/pkg"

	"github.com/sunshineplan/imgconv"
)

type convertTool struct {
	log pkg.Logger
}

func NewOptimizeTool(logger pkg.Logger) ConvertTool {
	var ConvertTool ConvertTool = &convertTool{
		log: logger,
	}
	return ConvertTool
}

type ConvertTool interface {
	ConvertImageToJpeg(input *image.Image) (output image.Image, err error)
	ConvertImageToByteArray(image *image.Image) (output *[]byte, err error)
}

func (c *convertTool) ConvertImageToJpeg(input *image.Image) (output image.Image, err error) {
	var buf bytes.Buffer
	err = imgconv.Write(&buf, *input, imgconv.FormatOption{Format: imgconv.JPEG})
	if err != nil {
		c.log.Errorf("converter can't convert image, reason: %v", err.Error())
		err = image_errors.ErrCantConvert
		return
	}
	img, err := jpeg.Decode(bytes.NewReader(buf.Bytes()))
	if err != nil {
		c.log.Errorf("converter can't decode image, reason: %v", err.Error())
		err = image_errors.ErrCantConvert
		return
	}
	output = img
	return
}

func (c *convertTool) ConvertImageToByteArray(image *image.Image) (output *[]byte, err error) {
	buf := new(bytes.Buffer)
	err = jpeg.Encode(buf, *image, nil)
	if err != nil {
		c.log.Errorf("converter can't convert image to byte, reason: %v", err.Error())
		err = image_errors.ErrCantConvert
		return
	}
	byteArray := buf.Bytes()
	output = &byteArray
	return
}
