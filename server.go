package main

import (
	"bufio"
	"fmt"
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

func printHtml(w http.ResponseWriter, r *http.Request) {
	html, err := readLines("./index.html")
	if err != nil {
		return
	}
	_, err = fmt.Fprint(w, html)
	if err != nil {
		return
	}
}

func main() {
	http.HandleFunc("/", printHtml)
	//http.HandleFunc("/login", printHtml)
	//http.HandleFunc("/signup", printHtml)

	fmt.Println("starting server at :8080")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatalf("server error %s", err)
	}
}
