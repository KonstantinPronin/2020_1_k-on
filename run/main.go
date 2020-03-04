package main

import (
	"2020_1_k-on/application/server"
	"log"
)

func main() {
	serv := server.NewServer(":8080")
	log.Println("starting server at :8080")
	log.Fatal(serv.ListenAndServe())
}
