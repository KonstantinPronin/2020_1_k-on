package repository

import (
	"2020_1_k-on/application/models"
	"2020_1_k-on/infrastructure"
	"bufio"
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
	"sync"
)

//Интерфейсы запросов к бд

type FilmsList struct {
	mutex sync.RWMutex
	films map[string]*models.Film
	count uint
}

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

func CreateFilmList() *FilmsList {
	film := []models.Film{}
	filmList := FilmsList{}
	filmList = FilmsList{
		mutex: sync.RWMutex{},
		films: make(map[string]*models.Film),
		count: 0,
	}
	str, _ := readLines(infrastructure.DBpath)
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
	filmList.mutex.Lock()
	defer filmList.mutex.Unlock()
	fArr := []models.Film{}
	for _, val := range filmList.films {
		fArr = append(fArr, *val)
	}
	js, err := json.Marshal(fArr)
	if err != nil {
		return false
	}
	err = ioutil.WriteFile(infrastructure.DBpath, js, os.ModePerm)
	if err != nil {
		return false
	}
	return true
}

func (filmList *FilmsList) Create(film *models.Film) (models.Film, bool) {
	filmList.mutex.Lock()
	id := filmList.count + 1
	film.ID = id
	_, ok := filmList.films[film.Name]
	if ok {
		return models.Film{}, false
	}
	filmList.films[film.Name] = film
	filmList.count++
	filmList.mutex.Unlock()
	filmList.UpdateFilmList()
	return *filmList.films[film.Name], true
}

func (filmList *FilmsList) GetById(id uint) (*models.Film, bool) {
	filmList.mutex.RLock()
	defer filmList.mutex.RUnlock()
	for _, f := range filmList.films {
		if f.ID == id {
			return f, true
		}
	}
	return nil, false
}

func (filmList *FilmsList) GetByName(name string) (*models.Film, bool) {
	filmList.mutex.RLock()
	defer filmList.mutex.RUnlock()

	f, ok := filmList.films[name]
	if !ok {
		return nil, false
	}
	return f, true
}
