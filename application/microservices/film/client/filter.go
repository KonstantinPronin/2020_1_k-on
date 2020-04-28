package client

import (
	"context"
	"github.com/go-park-mail-ru/2020_1_k-on/application/microservices/film/api"
	"github.com/go-park-mail-ru/2020_1_k-on/application/models"
	"go.uber.org/zap"
	"google.golang.org/grpc"
)

type FilmFilterClient struct {
	client api.FilmFilterClient
	gConn  *grpc.ClientConn
	logger *zap.Logger
}

func NewFilmFilterClient(host, port string, logger *zap.Logger) (*FilmFilterClient, error) {
	gConn, err := grpc.Dial(
		host+port,
		grpc.WithInsecure(),
	)
	if err != nil {
		return nil, err
	}

	return &FilmFilterClient{
		client: api.NewFilmFilterClient(gConn),
		gConn:  gConn,
		logger: logger,
	}, nil
}

func (f *FilmFilterClient) GetFilteredFilms(fields map[string][]string) (models.Films, bool) {
	films, err := f.client.GetFilteredFilms(context.Background(), convertToFilter(fields))
	if err != nil {
		f.logger.Error(err.Error())
		return nil, false
	}

	return convertFromFilms(films), true
}

func (f *FilmFilterClient) GetFilterFields() (map[string]models.Genres, bool) {
	filterFields, err := f.client.GetFilterFields(context.Background(), &api.EmptyArgs{})
	if err != nil {
		f.logger.Error(err.Error())
		return nil, false
	}

	return convertFromFilterFields(filterFields), true
}

func (f *FilmFilterClient) Close() {
	if err := f.gConn.Close(); err != nil {
		err = f.gConn.Close()
		f.logger.Error("error while closing grpc connection")
	}
}
