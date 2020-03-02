package main

import (
	"encoding/base64"
	"io/ioutil"
	"os"
	filepath2 "path/filepath"
)

const permissions = 0644

type IoFile interface {
	WriteFile(filename string, data []byte, perm os.FileMode) error
	ReadFile(filename string) ([]byte, error)
}

type IoUtil struct{}

func (i *IoUtil) WriteFile(filename string, data []byte, perm os.FileMode) error {
	return ioutil.WriteFile(filename, data, perm)
}

func (i *IoUtil) ReadFile(filename string) ([]byte, error) {
	return ioutil.ReadFile(filename)
}

type FileHandler struct {
	ioFile IoFile
}

func NewFileHandler(ioFile IoFile) *FileHandler {
	return &FileHandler{ioFile: ioFile}
}

func (fileHandler *FileHandler) WriteFile(directory string, filename string, data []byte) (string, error) {
	filepath := filepath2.Join(directory, filename)
	err := fileHandler.ioFile.WriteFile(filepath, data, permissions)
	if err != nil {
		return "", err
	}

	return filepath, nil
}

func (fileHandler *FileHandler) ReadFile(filepath string) (string, error) {
	data, err := fileHandler.ioFile.ReadFile(filepath)
	if err != nil {
		return "", err
	}

	return base64.StdEncoding.EncodeToString(data), nil
}
