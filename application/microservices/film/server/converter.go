package server

import (
	"github.com/go-park-mail-ru/2020_1_k-on/application/microservices/film/api"
	"github.com/go-park-mail-ru/2020_1_k-on/application/models"
)

func convertToFilterFields(fields map[string]models.Genres) *api.FilterFields {
	result := new(api.FilterFields)
	result.Fields = make(map[string]*api.Genres)
	for key, value := range fields {
		genres := new(api.Genres)
		for _, g := range value {
			genres.List = append(genres.List, &api.Genre{
				Name:      g.Name,
				Reference: g.Reference,
			})
		}
		result.Fields[key] = genres
	}
	return result
}

func convertToFilms(films models.Films) *api.Films {
	result := new(api.Films)
	for _, f := range films {
		result.List = append(result.List, &api.Film{
			Id:              uint64(f.ID),
			MainGenre:       f.MainGenre,
			RussianName:     f.RussianName,
			EnglishName:     f.EnglishName,
			TrailerLink:     f.TrailerLink,
			Rating:          f.Rating,
			ImdbRating:      f.ImdbRating,
			TotalVotes:      int64(f.TotalVotes),
			SumVotes:        int64(f.SumVotes),
			Description:     f.Description,
			Image:           f.Image,
			BackgroundImage: f.BackgroundImage,
			Country:         f.Country,
			Year:            int64(f.Year),
			AgeLimit:        int64(f.AgeLimit),
		})
	}
	return result
}

func convertFromFilter(filter *api.Filter) map[string][]string {
	result := make(map[string][]string)
	for key, value := range filter.Fields {
		result[key] = value.Value
	}
	return result
}
