package main

import (
	"bufio"
	"fmt"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"os"
)

func readLines(path string) (string, error) {
	file, err := os.Open(path)
	if err != nil {
		return "", err
	}
	defer file.Close()

	var lines string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines = lines + scanner.Text()
	}
	return lines, scanner.Err()
}

func main() {
	r := mux.NewRouter()

	filmHandler := createFilmHandler()
	r.HandleFunc("/films/create/", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case "GET":
			http.ServeFile(w, r, "../storage/index.html")
		case "POST":
			filmHandler.createFilm(w, r)
		}
	})

	r.HandleFunc("/films", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		filmHandler.getFilmsList(w, r)
	})

	r.HandleFunc("/films/{id:[0-9]+}/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		filmHandler.getFilm(w, r)
	})
	fmt.Println("starting server at :8080")
	err := http.ListenAndServe(":8080", r)
	if err != nil {
		log.Fatalf("server error %s", err)
	}
}
