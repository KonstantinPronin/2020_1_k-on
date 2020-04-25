package server

import (
	"context"
	"github.com/go-park-mail-ru/2020_1_k-on/application/microservices/film/api"
	"github.com/go-park-mail-ru/2020_1_k-on/application/microservices/film/filter"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type FilmFilter struct {
	usecase filter.UseCase
}

func NewFilmFilter(usecase filter.UseCase) *FilmFilter {
	return &FilmFilter{usecase: usecase}
}

func (f *FilmFilter) GetFilterFields(context.Context, *api.EmptyArgs) (*api.FilterFields, error) {
	fields, ok := f.usecase.FilterFilmData()
	if !ok {
		return nil, status.Error(codes.InvalidArgument, "")
	}

	return convertToFilterFields(fields), nil
}

func (f *FilmFilter) GetFilteredFilms(ctx context.Context, filter *api.Filter) (*api.Films, error) {
	arg := convertFromFilter(filter)
	films, ok := f.usecase.FilterFilmList(arg)
	if !ok {
		return nil, status.Error(codes.InvalidArgument, "")
	}

	return convertToFilms(films), nil
}
