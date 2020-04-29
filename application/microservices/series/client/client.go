package client

import "github.com/go-park-mail-ru/2020_1_k-on/application/models"

type ISeriesFilterClient interface {
	GetFilteredSeries(fields map[string][]string) (models.SeriesArr, bool)
	GetFilterFields() (map[string]models.Genres, bool)
}
