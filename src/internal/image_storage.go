package internal

import "image_storage/src/internal/domain"

type ImageStorage interface {
	Save(image *domain.MyImage) error
	Get(image *domain.MyImage) error
}
