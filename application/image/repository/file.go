package repository

import (
	"encoding/base64"
	"github.com/go-park-mail-ru/2020_1_k-on/application/image"
	"github.com/go-park-mail-ru/2020_1_k-on/pkg/constants"
	"github.com/google/uuid"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
)

type ImageRepository struct{}

func NewImageRepository() image.Repository {
	return &ImageRepository{}
}

func (rep *ImageRepository) Save(data io.Reader) (string, error) {
	filename := uuid.New().String()
	path := filepath.Join(constants.ImgDir, filename)

	dst, err := os.Create(path)
	if err != nil {
		return "", err
	}
	defer dst.Close()

	if _, err = io.Copy(dst, data); err != nil {
		return "", err
	}

	return filename, nil
}

func (rep *ImageRepository) Get(filename string) (string, error) {
	path := filepath.Join(constants.ImgDir, filename)
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return "", err
	}

	return base64.StdEncoding.EncodeToString(data), nil
}
