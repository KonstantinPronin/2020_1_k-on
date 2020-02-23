package main

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
)

func getFilmsList(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	film := []Film{}
	str, _ := readLines("../storage/films.json")
	data := []byte(str)
	err := json.Unmarshal(data, &film)

	if err != nil {
		return
	}
	for k := range film {
		_, err = fmt.Fprintf(w, "ID:%d, Name:%s, Year: %d \n", film[k].ID, film[k].Name, film[k].YearCreated)
	}
	if err != nil {
		return
	}
	return
}

func getFilm(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, `{"error":"bad id"}`, 400)
		return
	}
	film := []Film{}
	str, _ := readLines("../storage/films.json")
	data := []byte(str)
	err = json.Unmarshal(data, &film)
	if err != nil {
		return
	}
	for k := range film {
		if film[k].ID == id {
			_, err = fmt.Fprintf(w, "ID:%d, Name:%s, Year: %d \n", film[k].ID, film[k].Name, film[k].YearCreated)
			return
		}
	}
	http.Error(w, "not found", 404)
	return

}

func createFilm(w http.ResponseWriter, r *http.Request) {
	film := []Film{}
	str, _ := readLines("../storage/films.json")
	data := []byte(str)
	err := json.Unmarshal(data, &film)
	if err != nil {
		http.Error(w, "server error in unmarshal", 500)
		log.Fatalf("%s", err)
		return
	}

	if err := r.ParseForm(); err != nil {
		http.Error(w, "server error in parsing form", 500)
		log.Fatalf("%s", err)

		return
	}
	var id int = 0
	if len(film) > 0 {
		id = film[len(film)-1].ID + 1
	}
	name := r.FormValue("Name")
	year, err := strconv.Atoi(r.FormValue("Year"))
	if err != nil {
		fmt.Fprintf(w, "year must be int!: %v", err)
		return
	}
	film = append(film, Film{ID: id, Name: name, YearCreated: year})
	js, err := json.Marshal(film)
	if err != nil {
		http.Error(w, "server error in marshal", 500)
		log.Fatalf("%s", err)
		return
	}
	err = ioutil.WriteFile("../storage/films.json", js, os.ModePerm)
	if err != nil {
		http.Error(w, "server error in saving data", 500)
		log.Fatalf("%s", err)
	}
	http.Redirect(w, r, "/films/"+strconv.Itoa(id)+"/", 301)
}
