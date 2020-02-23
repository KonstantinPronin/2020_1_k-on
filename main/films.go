package main

type Film struct {
	ID          int
	Name        string
	YearCreated int
}

type FilmInput struct {
	Name        string `json:"name"`
	YearCreated int    `json:"year"`
}
