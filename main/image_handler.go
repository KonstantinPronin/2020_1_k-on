package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
)

const MaxFileSize int64 = 5 * 1024 * 1024
const ImageDirectory = "./media/img/"

type ImageJson struct {
	Image string `json:"image"`
}

type ImageHandler struct {
	fileHandler *FileHandler
}

func NewImageHandler(fileHandler *FileHandler) *ImageHandler {
	return &ImageHandler{fileHandler: fileHandler}
}

func (imageHandler *ImageHandler) AddImage(w http.ResponseWriter, r *http.Request) (string, bool) {
	err := r.ParseMultipartForm(MaxFileSize)
	if err != nil {
		http.Error(w, `{"error":"bad data form"}`, http.StatusBadRequest)
		return "", false
	}

	file, header, err := r.FormFile("file")
	if err != nil {
		http.Error(w, `{"error":"bad data form"}`, http.StatusBadRequest)
		return "", false
	}
	defer func() {
		err = file.Close()
		if err != nil {
			log.Println("Resource leaks")
		}
	}()

	data, err := ioutil.ReadAll(file)
	if err != nil {
		http.Error(w, `{"error":"cannot read file"}`, http.StatusInternalServerError)
		return "", false
	}

	filepath, err := imageHandler.fileHandler.WriteFile(ImageDirectory, header.Filename, data)
	if err != nil {
		http.Error(w, `{"error":"cannot save file"}`, http.StatusInternalServerError)
		return "", false
	}

	return filepath, true
}

func (imageHandler *ImageHandler) GetImage(w http.ResponseWriter, path string) {
	data, err := imageHandler.fileHandler.ReadFile(path)
	if err != nil {
		http.Error(w, `{"error":"`+string(err.Error())+`"}`, http.StatusInternalServerError)
		return
	}

	image := ImageJson{Image: data}
	err = json.NewEncoder(w).Encode(&image)
	if err != nil {
		http.Error(w, `{"error":"`+string(err.Error())+`"}`, http.StatusInternalServerError)
		return
	}
}
