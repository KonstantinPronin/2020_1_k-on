package main

import (
	"bytes"
	"encoding/json"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/require"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"sync"
	"testing"
)

const jsReq = `
{
	"ID":5,
	"name": "New",
	"AgeLimit": 3000,
	"Image": "path"
}
`

const jsResp = `[{"ID":10,"Name":"1","AgeLimit":1000,"Image":"path"}]`
const jsResp2 = `{"ID":10,"Name":"1","AgeLimit":1000,"Image":"path"}`

func TestCreateFilm(t *testing.T) {
	var film Film
	data := []byte(jsReq)
	json.Unmarshal(data, &film)
	body := bytes.NewReader([]byte(jsReq))
	r := httptest.NewRequest("POST", "/films/create", body)
	w := httptest.NewRecorder()
	handl := FilmHandler{films: CreateFilmList()}
	handl.createFilm(w, r)
	f, ok := handl.films.GetByName(film.Name)
	require.Equal(t, *f, film)
	require.True(t, ok)
	delete(handl.films.films, film.Name)
	handl.films.UpdateFilmList()
}

func TestGetFilmsList(t *testing.T) {
	film := []Film{
		{
			ID:       10,
			Name:     "1",
			AgeLimit: 1000,
			Image:    "path",
		},
	}
	filmList := FilmsList{
		mutex: sync.RWMutex{},
		films: make(map[string]*Film),
		count: 0,
	}
	for ind, val := range film {
		filmList.films[val.Name] = &(film[ind])
		filmList.count += 1
	}
	handler := FilmHandler{films: &filmList}

	r := httptest.NewRequest("GET", "/films", nil)
	w := httptest.NewRecorder()
	handler.getFilmsList(w, r)
	require.Equal(t, w.Code, http.StatusOK)
	data, _ := ioutil.ReadAll(w.Body)
	require.Equal(t, strings.Trim(string(data), "\n"), jsResp)
}

func TestGetFilm(t *testing.T) {
	film := []Film{
		{
			ID:       10,
			Name:     "1",
			AgeLimit: 1000,
			Image:    "path",
		},
	}
	filmList := FilmsList{
		mutex: sync.RWMutex{},
		films: make(map[string]*Film),
		count: 0,
	}
	for ind, val := range film {
		filmList.films[val.Name] = &(film[ind])
		filmList.count += 1
	}
	handler := FilmHandler{films: &filmList}

	r := httptest.NewRequest("GET", "/films/"+strconv.Itoa(int(film[0].ID))+"/", nil)
	r = mux.SetURLVars(r, map[string]string{
		"id": "10",
	})
	w := httptest.NewRecorder()
	handler.getFilm(w, r)
	require.Equal(t, w.Code, http.StatusOK)
	data, _ := ioutil.ReadAll(w.Body)
	require.Equal(t, strings.Trim(string(data), "\n"), jsResp2)
}
