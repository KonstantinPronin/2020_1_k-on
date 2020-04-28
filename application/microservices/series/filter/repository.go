package filter

import "github.com/go-park-mail-ru/2020_1_k-on/application/models"

type Repository interface {
	FilterSeriesList(fields map[string][]string) (*models.SeriesArr, bool)
	FilterSeriesData() (map[string]models.Genres, bool)
}
