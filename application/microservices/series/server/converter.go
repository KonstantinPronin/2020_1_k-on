package server

import (
	api "github.com/go-park-mail-ru/2020_1_k-on/application/microservices/series/api"
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

func convertToSeriesList(slist models.SeriesArr) *api.SeriesList {
	result := new(api.SeriesList)
	for _, s := range slist {
		result.List = append(result.List, &api.Series{
			Id:              uint64(s.ID),
			MainGenre:       s.MainGenre,
			RussianName:     s.RussianName,
			EnglishName:     s.EnglishName,
			TrailerLink:     s.TrailerLink,
			Rating:          s.Rating,
			ImdbRating:      s.ImdbRating,
			TotalVotes:      int64(s.TotalVotes),
			SumVotes:        int64(s.SumVotes),
			Description:     s.Description,
			Image:           s.Image,
			BackgroundImage: s.BackgroundImage,
			Country:         s.Country,
			YearFirst:       int64(s.YearFirst),
			YearLast:        int64(s.YearLast),
			AgeLimit:        int64(s.AgeLimit),
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
