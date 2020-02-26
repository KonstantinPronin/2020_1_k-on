package main

import (
	"encoding/base64"
	"io/ioutil"
)

const ImageDirectory string = string("./media/img/")
const MaxFileSize int64 = 1000000

type ImageJson struct {
	Image string `json:"image"`
}

func WriteFile(directory string, filename string, data []byte) (string, error) {
	filepath := directory + filename
	err := ioutil.WriteFile(filepath, data, 0644)
	if err != nil {
		return "", err
	}

	return filepath, nil
}

func ReadFile(filepath string) (string, error) {
	data, err := ioutil.ReadFile(filepath)
	if err != nil {
		return "", err
	}

	return base64.StdEncoding.EncodeToString(data), nil
}
