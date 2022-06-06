package mock

import "image"

type ConverterToolMock struct {
	FuncConvertImageToJpeg      func(input *image.Image) (output image.Image, err error)
	FuncConvertImageToByteArray func(myImg *image.Image) (output *[]byte, err error)
}

func (c *ConverterToolMock) ConvertImageToJpeg(input *image.Image) (output image.Image, err error) {
	return c.FuncConvertImageToJpeg(input)
}

func (c *ConverterToolMock) ConvertImageToByteArray(image *image.Image) (output *[]byte, err error) {
	return c.FuncConvertImageToByteArray(image)
}
