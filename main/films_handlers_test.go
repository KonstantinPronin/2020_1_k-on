package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/stretchr/testify/require"
	"net/http/httptest"
	"testing"
)

const jsonIn = `[
  {
    "ID": 6,
    "Name": "New",
    "YearCreated": 3000
  }
]`

const jsReq = `{
  "name": "New",
  "YearCreated": 3000
}`

func TestCreateFilm(t *testing.T) {
	film := []Film{}
	data := []byte(jsonIn)
	json.Unmarshal(data, &film)
	fmt.Print(film)
	body := bytes.NewReader([]byte(jsReq))
	r := httptest.NewRequest("POST", "/films/create/", body)
	w := httptest.NewRecorder()
	handl := createFilmHandler()
	handl.createFilm(w, r)
	f, ok := handl.films.GetByName("New")
	require.Equal(t, *f, film[0])
	require.True(t, ok)
	delete(handl.films.films, "New")
	handl.films.UpdateFilmList()
}
