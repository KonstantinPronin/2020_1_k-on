package http

import (
	"2020_1_k-on/application/film"
	"2020_1_k-on/application/models"
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"strconv"
)

type FilmHandler struct {
	usecase film.Usecase
}

func NewFilmHandler(router *mux.Router, usecase film.Usecase) {
	handler := &FilmHandler{
		usecase: usecase,
	}
	router.HandleFunc("/films/{id:[0-9]+}", handler.GetFilm).Methods(http.MethodGet)
	router.HandleFunc("/films", handler.CreateFilm).Methods(http.MethodPost)
	router.HandleFunc("/films", handler.GetFilmList).Methods(http.MethodGet)

}

func (fh FilmHandler) GetFilm(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		fmt.Fprint(w, `{"Answer":"BAD"}`)
		return
	}
	f, ok := fh.usecase.GetFilm(uint(id))
	if !ok {
		fmt.Fprint(w, `{"Answer":"BAD"}`)
		return
	}
	json.NewEncoder(w).Encode(f)
}

func (fh FilmHandler) GetFilmList(w http.ResponseWriter, r *http.Request) {
	f := fh.usecase.GetFilmsList()
	json.NewEncoder(w).Encode(f)
}

func (fh FilmHandler) CreateFilm(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	decoder := json.NewDecoder(r.Body)
	var film models.Film
	err := decoder.Decode(&film)
	if err != nil {
		http.Error(w, `{"error":"can't parse json'"}`, 500)
		log.Printf("%s", err)
		return
	}

	f := fh.usecase.CreateFilm(film)
	//if ok {
	//	json.NewEncoder(w).Encode(film)
	//} else {
	//	http.Error(w, "can't update database", http.StatusInternalServerError)
	//}

	json.NewEncoder(w).Encode(f)
}
