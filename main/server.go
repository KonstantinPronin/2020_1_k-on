package main

import (
	"bufio"
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
	router := mux.NewRouter()

	filmHandler := createFilmHandler()
	userHandler := createUserHandler()

	router.HandleFunc("/films/create/", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case "GET":
			http.ServeFile(w, r, "../storage/index.html")
		case "POST":
			filmHandler.createFilm(w, r)
		}
	})

	router.HandleFunc("/films", filmHandler.getFilmsList)
	router.HandleFunc("/films/{id:[0-9]+}/", filmHandler.getFilm)
	router.HandleFunc("/login", userHandler.Login)
	router.HandleFunc("/logout", userHandler.Logout)
	router.HandleFunc("/signup", userHandler.Add)
	router.HandleFunc("/user", userHandler.Profile)
	router.HandleFunc("/user/image", userHandler.UploadImage)
	router.HandleFunc("/user/{id:[0-9]+}/image", userHandler.GetImage)

	log.Println("starting server at :8080")
	err := http.ListenAndServe(":8080", r)
	if err != nil {
		log.Fatalf("server error %s", err)
	}
}
