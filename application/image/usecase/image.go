package usecase

import (
	"github.com/go-park-mail-ru/2020_1_k-on/application/image"
	"go.uber.org/zap"
	"io"
)

type Image struct {
	images image.Repository
	logger *zap.Logger
}

func NewImage(images image.Repository, logger *zap.Logger) image.UseCase {
	return &Image{
		images: images,
		logger: logger,
	}
}

func (i *Image) Get(filepath string) (string, error) {
	return i.images.Get(filepath)
}

func (i *Image) Save(image io.Reader) (string, error) {
	return i.images.Save(image)
}
