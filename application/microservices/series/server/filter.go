package server

import (
	"context"
	api "github.com/go-park-mail-ru/2020_1_k-on/application/microservices/series/api"
	"github.com/go-park-mail-ru/2020_1_k-on/application/microservices/series/filter"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type SeriesFilter struct {
	usecase filter.UseCase
}

func NewSeriesFilter(usecase filter.UseCase) *SeriesFilter {
	return &SeriesFilter{usecase: usecase}
}

func (f *SeriesFilter) GetFilterFields(context.Context, *api.EmptyArgs) (*api.FilterFields, error) {
	fields, ok := f.usecase.FilterSeriesData()
	if !ok {
		return nil, status.Error(codes.InvalidArgument, "")
	}

	return convertToFilterFields(fields), nil
}

func (f *SeriesFilter) GetFilteredSeries(ctx context.Context, filter *api.Filter) (*api.SeriesList, error) {
	arg := convertFromFilter(filter)
	slist, ok := f.usecase.FilterSeriesList(arg)
	if !ok {
		return nil, status.Error(codes.InvalidArgument, "")
	}

	return convertToSeriesList(*slist), nil
}
