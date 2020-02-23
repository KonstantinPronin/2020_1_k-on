package main

import (
	"encoding/json"
	"github.com/stretchr/testify/require"
	"io/ioutil"
	"net/http/httptest"
	"net/url"
	"os"
	"strings"
	"testing"
)

func TestCreateFilm(t *testing.T) {
	form := url.Values{}
	form.Add("Name", "Abc")
	form.Add("Year", "2000")

	r := httptest.NewRequest("POST", "/films/create/", strings.NewReader(form.Encode()))
	r.Form = form
	w := httptest.NewRecorder()
	createFilm(w, r)
	//if w.Code != http.StatusOK {
	//	t.Error("status is not ok")
	//}
	film := []Film{}
	str, _ := readLines("../storage/films.json")
	data := []byte(str)
	err := json.Unmarshal(data, &film)
	require.NoError(t, err)
	OkFilm := Film{
		ID:          len(film),
		Name:        "Abc",
		YearCreated: 2000,
	}
	require.Equal(t, OkFilm, film[len(film)-1])

	js, _ := json.Marshal(film[:len(film)-1])
	ioutil.WriteFile("../storage/films.json", js, os.ModePerm)
}
