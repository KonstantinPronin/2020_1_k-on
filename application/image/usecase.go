package image

import (
	"io"
)

type UseCase interface {
	Get(filepath string) (string, error)
	Save(image io.Reader) (string, error)
}
