package main

import (
	"github.com/gorilla/mux"
	"log"
	"net/http"
)

func main() {
	router := mux.NewRouter()

	filmHandler := createFilmHandler()
	userHandler := NewUserHandler()

	router.HandleFunc("/films", filmHandler.createFilm).Methods(http.MethodOptions, http.MethodPost)
	router.HandleFunc("/films", filmHandler.getFilmsList).Methods(http.MethodOptions, http.MethodGet)
	router.HandleFunc("/films/{id:[0-9]+}", filmHandler.getFilm).Methods(http.MethodOptions, http.MethodGet)
	router.HandleFunc("/login", userHandler.Login)
	router.HandleFunc("/logout", userHandler.Logout)
	router.HandleFunc("/signup", userHandler.Add)
	router.HandleFunc("/user", userHandler.Get).Methods(http.MethodGet, http.MethodOptions)
	router.HandleFunc("/user", userHandler.Update).Methods(http.MethodPost, http.MethodOptions)
	router.HandleFunc("/user/image", userHandler.UploadImage)
	router.HandleFunc("/user/{id:[0-9]+}/image", userHandler.GetImage)
	router.HandleFunc("/film/image", filmHandler.UploadImageFilm)
	router.HandleFunc("/film/{id:[0-9]+}/image", filmHandler.GetImageFilm)

	log.Println("starting server at :8080")
	log.Fatal(http.ListenAndServe(":8080", router))
}
