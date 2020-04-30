package client

import (
	"context"
	api "github.com/go-park-mail-ru/2020_1_k-on/application/microservices/series/api"
	"github.com/go-park-mail-ru/2020_1_k-on/application/models"
	traceutils "github.com/opentracing-contrib/go-grpc"
	"github.com/opentracing/opentracing-go"
	"go.uber.org/zap"
	"google.golang.org/grpc"
)

type SeriesFilterClient struct {
	client api.SeriesFilterClient
	gConn  *grpc.ClientConn
	logger *zap.Logger
}

func NewSeriesFilterClient(host, port string, logger *zap.Logger, tracer opentracing.Tracer) (*SeriesFilterClient, error) {
	gConn, err := grpc.Dial(
		host+port,
		grpc.WithInsecure(),
		grpc.WithUnaryInterceptor(traceutils.OpenTracingClientInterceptor(tracer)),
	)
	if err != nil {
		return nil, err
	}

	return &SeriesFilterClient{
		client: api.NewSeriesFilterClient(gConn),
		gConn:  gConn,
		logger: logger,
	}, nil
}

func (s *SeriesFilterClient) GetFilteredSeries(fields map[string][]string) (models.SeriesArr, bool) {
	slist, err := s.client.GetFilteredSeries(context.Background(), convertToFilter(fields))
	if err != nil {
		s.logger.Error(err.Error())
		return nil, false
	}

	return convertFromSeriesList(slist), true
}

func (s *SeriesFilterClient) GetFilterFields() (map[string]models.Genres, bool) {
	filterFields, err := s.client.GetFilterFields(context.Background(), &api.EmptyArgs{})
	if err != nil {
		s.logger.Error(err.Error())
		return nil, false
	}

	return convertFromFilterFields(filterFields), true
}

func (s *SeriesFilterClient) Close() {
	if err := s.gConn.Close(); err != nil {
		err = s.gConn.Close()
		s.logger.Error("error while closing grpc connection")
	}
}
