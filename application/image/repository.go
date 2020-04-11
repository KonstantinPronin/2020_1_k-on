package image

import "io"

type Repository interface {
	Get(filepath string) (string, error)
	Save(image io.Reader) (string, error)
}
