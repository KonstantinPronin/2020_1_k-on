package main

import (
	"encoding/json"
	"github.com/stretchr/testify/require"
	"sync"
	"testing"
)

const FilmsData = `[
{
	"ID": 1,
	"Name": "Avatar",
	"YearCreated": 2009
},
{
"ID": 2,
"Name": "Avengers",
"YearCreated": 2012
}]`

const OkData = `Avatar`
const FailData = `Aang`
const OkYear = 2000

func prepare(str string) FilmsList {
	film := []Film{}
	data := []byte(str)
	json.Unmarshal(data, &film)
	filmList := FilmsList{
		mutex: sync.RWMutex{},
		films: make(map[string]*Film),
		count: 0,
	}
	for ind, val := range film {
		filmList.films[val.Name] = &(film[ind])
		filmList.count += 1
	}
	return filmList
}

func TestFilmsList_Contains(t *testing.T) {
	fl := prepare(FilmsData)
	require.True(t, fl.Contains(OkData))
	require.False(t, fl.Contains(FailData))
}

func TestFilmsList_GetById(t *testing.T) {
	fl := prepare(FilmsData)
	f, ok := fl.GetById(1)
	require.Equal(t, f.Name, OkData)
	require.NotEqual(t, f.Name, FailData)
	require.True(t, ok)
}

func TestFilmsList_GetByName(t *testing.T) {
	fl := prepare(FilmsData)
	f, ok := fl.GetByName(OkData)
	require.Equal(t, f.Name, OkData)
	require.NotEqual(t, f.Name, FailData)
	require.True(t, ok)
}

func TestFilmsList_Add(t *testing.T) {
	fl := prepare(FilmsData)
	f1 := Film{Name: OkData, AgeLimit: OkYear}
	f2 := Film{Name: FailData, AgeLimit: OkYear}
	fl.Add(&f1)
	fl.Add(&f2)
	require.Equal(t, fl.films[OkData], &f1)
	require.NotEqual(t, fl.films[OkData], &f2)
}
