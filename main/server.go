package main

import (
	"github.com/gorilla/mux"
	"log"
	"net/http"
)

func main() {
	router := mux.NewRouter()

	filmHandler := createFilmHandler()
	userHandler := createUserHandler()

	router.HandleFunc("/films/create/", filmHandler.createFilm)
	router.HandleFunc("/films", filmHandler.getFilmsList)
	router.HandleFunc("/films/{id:[0-9]+}/", filmHandler.getFilm)
	router.HandleFunc("/login", userHandler.Login)
	router.HandleFunc("/logout", userHandler.Logout)
	router.HandleFunc("/signup", userHandler.Add)
	router.HandleFunc("/user", userHandler.Get).Methods(http.MethodGet, http.MethodOptions)
	router.HandleFunc("/user", userHandler.Update).Methods(http.MethodPost, http.MethodOptions)
	router.HandleFunc("/user/image", userHandler.UploadImage)
	router.HandleFunc("/user/{id:[0-9]+}/image", userHandler.GetImage)

	log.Println("starting server at :8080")
	log.Fatal(http.ListenAndServe(":8080", router))
}
