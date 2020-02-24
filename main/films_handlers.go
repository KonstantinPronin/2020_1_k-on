package main

import (
	"fmt"
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
	for _, val := range filmHandler.films.films {
		fmt.Fprintf(w, "%d ", val.ID)
		fmt.Fprintf(w, "%s ", val.Name)
		fmt.Fprintf(w, "%d \n", val.YearCreated)

	}
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
	f, ok := filmHandler.films.GetById(uint(id))
	if !ok {
		http.Error(w, `{"error":"not found this film"}`, 404)
		return
	} else {
		fmt.Fprintf(w, "%d ", f.ID)
		fmt.Fprintf(w, "%s ", f.Name)
		fmt.Fprintf(w, "%d \n", f.YearCreated)
	}
}

func (filmHandler *FilmHandler) createFilm(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		http.Error(w, "server error in parsing form", 500)
		log.Fatalf("%s", err)

		return
	}
	name := r.FormValue("Name")
	_, is := filmHandler.films.GetByName(name)
	if is {
		fmt.Fprint(w, "film with this name already exists:")
		return
	}
	year, err := strconv.Atoi(r.FormValue("Year"))
	if err != nil {
		fmt.Fprintf(w, "year must be int!: %s", err)
		return
	}
	film := Film{
		ID:          filmHandler.films.count + 1,
		Name:        name,
		YearCreated: year,
	}

	filmHandler.films.Add(&film)
	ok := filmHandler.films.UpdateFilmList()
	if ok {
		http.Redirect(w, r, "/films/"+strconv.Itoa(int(film.ID))+"/", 301)
	} else {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
