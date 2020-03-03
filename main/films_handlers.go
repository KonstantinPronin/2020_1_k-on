package main

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"strconv"
)

type FilmHandler struct {
	films        *FilmsList
	imageHandler *ImageHandler
}

func createFilmHandler() *FilmHandler {
	return &FilmHandler{
		films:        CreateFilmList(),
		imageHandler: NewImageHandler(NewFileHandler(new(IoUtil))),
	}
}

func (filmHandler *FilmHandler) getFilmsList(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	w.Header().Set("Access-Control-Allow-Methods", "GET")
	w.Header().Set("Access-Control-Allow-Origin", r.Header.Get("Origin"))
	w.Header().Set("Access-Control-Allow-Credentials", "true")

	if r.Method == "GET" {
		w.Header().Set("Content-Type", "application/json")
		var f []Film
		for _, val := range filmHandler.films.films {
			if val.Image != "" {
				data, err := filmHandler.imageHandler.fileHandler.ReadFile(val.Image)
				if err != nil {
					http.Error(w, `{"error":"`+string(err.Error())+`"}`, http.StatusInternalServerError)
					return
				}
				val.ImageBase64 = data
			}
			f = append(f, *val)
		}
		json.NewEncoder(w).Encode(f)
		return
	}
}

func (filmHandler *FilmHandler) getFilm(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	w.Header().Set("Access-Control-Allow-Methods", "GET")
	w.Header().Set("Access-Control-Allow-Origin", r.Header.Get("Origin"))
	w.Header().Set("Access-Control-Allow-Credentials", "true")

	if r.Method == "GET" {
		w.Header().Set("Content-Type", "application/json")
		vars := mux.Vars(r)
		id, err := strconv.Atoi(vars["id"])
		if err != nil {
			http.Error(w, `{"error":"bad id"}`, 400)
			return
		}
		film, ok := filmHandler.films.GetById(uint(id))
		if !ok {
			http.Error(w, `{"error":"not found this film"}`, 404)
			return
		} else {
			json.NewEncoder(w).Encode(film)
		}
		if film.Image != "" {
			data, err := filmHandler.imageHandler.fileHandler.ReadFile(film.Image)
			if err != nil {
				http.Error(w, `{"error":"`+string(err.Error())+`"}`, http.StatusInternalServerError)
				return
			}
			film.ImageBase64 = data
		}
	}
}

func (filmHandler *FilmHandler) createFilm(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	w.Header().Set("Access-Control-Allow-Methods", "POST")
	w.Header().Set("Access-Control-Allow-Origin", r.Header.Get("Origin"))
	w.Header().Set("Access-Control-Allow-Credentials", "true")

	if r.Method == "POST" {
		defer r.Body.Close()
		decoder := json.NewDecoder(r.Body)
		var film Film
		err := decoder.Decode(&film)
		if err != nil {
			http.Error(w, `{"error":"can't parse json'"}`, 500)
			log.Printf("%s", err)
			return
		}
		_, is := filmHandler.films.GetByName(film.Name)
		if is {
			http.Error(w, `bad film name`, http.StatusBadRequest)
			log.Print(w, "film with this name already exists:")
			return
		}
		filmHandler.films.Add(&film)
		ok := filmHandler.films.UpdateFilmList()
		if ok {
			json.NewEncoder(w).Encode(film)
		} else {
			http.Error(w, "can't update database", http.StatusInternalServerError)
		}
	}
}

func (filmHandler *FilmHandler) UploadImageFilm(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	w.Header().Set("Access-Control-Allow-Methods", "PUT")
	w.Header().Set("Access-Control-Allow-Origin", r.Header.Get("Origin"))
	w.Header().Set("Access-Control-Allow-Credentials", "true")

	if r.Method != http.MethodPut {
		return
	}
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, `{"error":"bad id"}`, 400)
		return
	}
	film, ok := filmHandler.films.GetById(uint(id))
	if !ok {
		http.Error(w, `{"error":"no film"}`, http.StatusNotFound)
		return
	}

	film.Image, ok = filmHandler.imageHandler.AddImage(w, r)
	if !ok {
		return
	}
	filmHandler.films.films[film.Name].Image = film.Image
	ok = filmHandler.films.UpdateFilmList()
	if ok {
		json.NewEncoder(w).Encode(film)
	} else {
		http.Error(w, `{"error":"can't update database"}`, http.StatusInternalServerError)
	}
}

func (filmHandler *FilmHandler) GetImageFilm(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	w.Header().Set("Access-Control-Allow-Methods", "GET")
	w.Header().Set("Access-Control-Allow-Origin", r.Header.Get("Origin"))
	w.Header().Set("Access-Control-Allow-Credentials", "true")

	if r.Method != http.MethodGet {
		return
	}

	w.Header().Set("Content-Type", "application/json")

	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if id < 0 || err != nil {
		http.Error(w, `{"error":"bad id"}`, http.StatusBadRequest)
		return
	}

	film, ok := filmHandler.films.GetById(uint(id))
	if !ok {
		http.Error(w, `{"error":"no user"}`, http.StatusNotFound)
		return
	}

	filmHandler.imageHandler.GetImage(w, film.Image)
}
