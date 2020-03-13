package main

import (
	"2020_1_k-on/application/server"
	"2020_1_k-on/infrastructure"
	"log"
)

func main() {
	conn, err := infrastructure.InitDatabaseConnection()
	if err != nil {
		log.Print(err)
	}
	serv := server.NewServer(":8080", conn)
	log.Println("starting server at :8080")
	log.Fatal(serv.ListenAndServe())
}
