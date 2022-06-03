package tools

import "image_storage/src/internal/domain"

type OptimizeTool struct{}

func NewOptimizeTool() *OptimizeTool {
	return &OptimizeTool{}
}

func (o *OptimizeTool) Optimize(input *domain.MyImage, quality string) (output *domain.MyImage, err error) {
	//bimg logic
	output = new(domain.MyImage)
	output.Id = input.Id
	output.Quality = quality
	output.Data = input.Data
	return
}
