package main

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"strconv"
)

type FilmHandler struct {
	films *FilmsList
}

func createFilmHandler() *FilmHandler {
	return &FilmHandler{
		films: CreateFilmList(),
	}
}

func (filmHandler *FilmHandler) getFilmsList(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var f []Film
	for _, val := range filmHandler.films.films {
		f = append(f, *val)
	}
	json.NewEncoder(w).Encode(f)
	return
}

func (filmHandler *FilmHandler) getFilm(w http.ResponseWriter, r *http.Request) {
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
}

func (filmHandler *FilmHandler) createFilm(w http.ResponseWriter, r *http.Request) {
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
		//http.Redirect(w, r, "/films/"+strconv.Itoa(int(film.ID))+"/", 301)
		json.NewEncoder(w).Encode(film)
	} else {
		http.Error(w, "can't update database", http.StatusInternalServerError)
	}
}
