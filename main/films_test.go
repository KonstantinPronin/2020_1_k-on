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

func TestFilmsList_Contains(t *testing.T) {
	film := []Film{}
	data := []byte(FilmsData)
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
	require.True(t, filmList.Contains(OkData))
	require.False(t, filmList.Contains(FailData))
}

func TestFilmsList_GetById(t *testing.T) {
	film := []Film{}
	data := []byte(FilmsData)
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
	f, ok := filmList.GetById(1)
	require.Equal(t, f.Name, OkData)
	require.NotEqual(t, f.Name, FailData)
	require.True(t, ok)
}

func TestFilmsList_GetByName(t *testing.T) {
	film := []Film{}
	data := []byte(FilmsData)
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
	f, ok := filmList.GetByName(OkData)
	require.Equal(t, f.Name, OkData)
	require.NotEqual(t, f.Name, FailData)
	require.True(t, ok)
}

func TestFilmsList_Add(t *testing.T) {
	film := []Film{}
	data := []byte(FilmsData)
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
	f1 := Film{Name: OkData, YearCreated: OkYear}
	f2 := Film{Name: FailData, YearCreated: OkYear}
	filmList.Add(&f1)
	filmList.Add(&f2)
	require.Equal(t, filmList.films[OkData], &f1)
	require.NotEqual(t, filmList.films[OkData], &f2)
}
