package http

import (
	"2020_1_k-on/application/film"
	"2020_1_k-on/application/models"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/mailru/easyjson"
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
		resp, _ := models.Generate(400, `{"Answer":"Bad id param"}`).MarshalJSON()
		w.Write(resp)
		return
	}
	f, ok := fh.usecase.GetFilm(uint(id))
	if !ok {
		resp, _ := models.Generate(404, `{"Answer":"Not found"}`).MarshalJSON()
		w.Write(resp)
		return
	}
	resp, _ := models.Generate(200, f).MarshalJSON()
	w.Write(resp)
}

func (fh FilmHandler) GetFilmList(w http.ResponseWriter, r *http.Request) {
	f := fh.usecase.GetFilmsList()
	resp, _ := models.Generate(200, f).MarshalJSON()
	w.Write(resp)
}

func (fh FilmHandler) CreateFilm(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	var film models.Film
	err := easyjson.UnmarshalFromReader(r.Body, &film)
	if err != nil {
		resp, _ := models.Generate(500, `{"error":"can't parse json'"}`).MarshalJSON()
		w.Write(resp)
		return
	}
	f, ok := fh.usecase.CreateFilm(film)
	if !ok {
		resp, _ := models.Generate(500, `{"error":"can't update db'"}`).MarshalJSON()
		w.Write(resp)
		return
	}
	resp, _ := models.Generate(200, f).MarshalJSON()
	_, err = w.Write(resp)
	if err != nil {
		fmt.Print(err)
		return
	}
}
