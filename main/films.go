package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
	"sync"
)

type Film struct {
	ID          uint
	Name        string
	YearCreated int
}

type FilmsList struct {
	mutex sync.RWMutex
	films map[string]*Film
	count uint
}

func CreateFilmList() *FilmsList {
	film := []Film{}
	filmList := FilmsList{}
	filmList = FilmsList{
		mutex: sync.RWMutex{},
		films: make(map[string]*Film),
		count: 0,
	}
	str, _ := readLines("../storage/films.json")
	data := []byte(str)
	err := json.Unmarshal(data, &film)
	if err != nil {
		log.Fatal(err)
	}

	for ind, val := range film {
		filmList.films[val.Name] = &(film[ind])
		filmList.count += 1
	}
	return &filmList
}

func (filmList *FilmsList) UpdateFilmList() bool {
	fArr := []Film{}
	for _, val := range filmList.films {
		fArr = append(fArr, *val)
	}
	js, err := json.Marshal(fArr)
	if err != nil {
		return false
	}
	err = ioutil.WriteFile("../storage/films.json", js, os.ModePerm)
	if err != nil {
		return false
	}
	return true
}

func (filmList *FilmsList) Contains(name string) bool {
	_, ok := filmList.films[name]
	return ok
}

func (filmList *FilmsList) Add(film *Film) uint {
	filmList.mutex.Lock()
	defer filmList.mutex.Unlock()
	id := filmList.count + 1
	film.ID = id
	filmList.films[film.Name] = film
	filmList.count++

	return id
}

func (filmList *FilmsList) GetById(id uint) (*Film, bool) {
	filmList.mutex.RLock()
	defer filmList.mutex.RUnlock()
	for _, film := range filmList.films {
		if film.ID == id {
			return film, true
		}
	}
	return nil, false
}

func (filmList *FilmsList) GetByName(name string) (*Film, bool) {
	filmList.mutex.RLock()
	defer filmList.mutex.RUnlock()

	if !filmList.Contains(name) {
		return nil, false
	}

	return filmList.films[name], true
}
