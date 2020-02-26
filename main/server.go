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

	router.HandleFunc("/films/create/", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case "GET":
			http.ServeFile(w, r, "../storage/index.html")
		case "POST":
			filmHandler.createFilm(w, r)
		}
	})

	router.HandleFunc("*", filmHandler.Cors).Methods("OPTIONS")

	router.HandleFunc("/films", filmHandler.getFilmsList)
	router.HandleFunc("/films/{id:[0-9]+}/", filmHandler.getFilm)
	router.HandleFunc("/login", userHandler.Login)
	router.HandleFunc("/logout", userHandler.Logout)
	router.HandleFunc("/signup", userHandler.Add)
	router.HandleFunc("/user", userHandler.Profile)
	router.HandleFunc("/user/image", userHandler.UploadImage)
	router.HandleFunc("/user/{id:[0-9]+}/image", userHandler.GetImage)

	log.Println("starting server at :8080")
	err := http.ListenAndServe(":8080", router)
	if err != nil {
		log.Fatalf("server error %s", err)
	}
}
