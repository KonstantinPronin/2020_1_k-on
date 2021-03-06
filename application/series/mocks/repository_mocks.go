// Code generated by MockGen. DO NOT EDIT.
// Source: ./application/series/repository.go

// Package mock_series is a generated GoMock package.
package mock_series

import (
	models "github.com/go-park-mail-ru/2020_1_k-on/application/models"
	gomock "github.com/golang/mock/gomock"
	reflect "reflect"
)

// MockRepository is a mock of Repository interface
type MockRepository struct {
	ctrl     *gomock.Controller
	recorder *MockRepositoryMockRecorder
}

// MockRepositoryMockRecorder is the mock recorder for MockRepository
type MockRepositoryMockRecorder struct {
	mock *MockRepository
}

// NewMockRepository creates a new mock instance
func NewMockRepository(ctrl *gomock.Controller) *MockRepository {
	mock := &MockRepository{ctrl: ctrl}
	mock.recorder = &MockRepositoryMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockRepository) EXPECT() *MockRepositoryMockRecorder {
	return m.recorder
}

// GetSeriesByID mocks base method
func (m *MockRepository) GetSeriesByID(id uint) (models.Series, bool) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetSeriesByID", id)
	ret0, _ := ret[0].(models.Series)
	ret1, _ := ret[1].(bool)
	return ret0, ret1
}

// GetSeriesByID indicates an expected call of GetSeriesByID
func (mr *MockRepositoryMockRecorder) GetSeriesByID(id interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetSeriesByID", reflect.TypeOf((*MockRepository)(nil).GetSeriesByID), id)
}

// GetSeriesSeasons mocks base method
func (m *MockRepository) GetSeriesSeasons(id uint) (models.Seasons, bool) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetSeriesSeasons", id)
	ret0, _ := ret[0].(models.Seasons)
	ret1, _ := ret[1].(bool)
	return ret0, ret1
}

// GetSeriesSeasons indicates an expected call of GetSeriesSeasons
func (mr *MockRepositoryMockRecorder) GetSeriesSeasons(id interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetSeriesSeasons", reflect.TypeOf((*MockRepository)(nil).GetSeriesSeasons), id)
}

// GetSeasonEpisodes mocks base method
func (m *MockRepository) GetSeasonEpisodes(id uint) (models.Episodes, bool) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetSeasonEpisodes", id)
	ret0, _ := ret[0].(models.Episodes)
	ret1, _ := ret[1].(bool)
	return ret0, ret1
}

// GetSeasonEpisodes indicates an expected call of GetSeasonEpisodes
func (mr *MockRepositoryMockRecorder) GetSeasonEpisodes(id interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetSeasonEpisodes", reflect.TypeOf((*MockRepository)(nil).GetSeasonEpisodes), id)
}

// GetSeriesGenres mocks base method
func (m *MockRepository) GetSeriesGenres(fid uint) (models.Genres, bool) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetSeriesGenres", fid)
	ret0, _ := ret[0].(models.Genres)
	ret1, _ := ret[1].(bool)
	return ret0, ret1
}

// GetSeriesGenres indicates an expected call of GetSeriesGenres
func (mr *MockRepositoryMockRecorder) GetSeriesGenres(fid interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetSeriesGenres", reflect.TypeOf((*MockRepository)(nil).GetSeriesGenres), fid)
}

// Search mocks base method
func (m *MockRepository) Search(word string, begin, end int) (models.SeriesArr, bool) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Search", word, begin, end)
	ret0, _ := ret[0].(models.SeriesArr)
	ret1, _ := ret[1].(bool)
	return ret0, ret1
}

// Search indicates an expected call of Search
func (mr *MockRepositoryMockRecorder) Search(word, begin, end interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Search", reflect.TypeOf((*MockRepository)(nil).Search), word, begin, end)
}

// GetSimilarFilms mocks base method
func (m *MockRepository) GetSimilarFilms(sid uint) (models.Films, bool) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetSimilarFilms", sid)
	ret0, _ := ret[0].(models.Films)
	ret1, _ := ret[1].(bool)
	return ret0, ret1
}

// GetSimilarFilms indicates an expected call of GetSimilarFilms
func (mr *MockRepositoryMockRecorder) GetSimilarFilms(sid interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetSimilarFilms", reflect.TypeOf((*MockRepository)(nil).GetSimilarFilms), sid)
}

// GetSimilarSeries mocks base method
func (m *MockRepository) GetSimilarSeries(sid uint) (models.SeriesArr, bool) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetSimilarSeries", sid)
	ret0, _ := ret[0].(models.SeriesArr)
	ret1, _ := ret[1].(bool)
	return ret0, ret1
}

// GetSimilarSeries indicates an expected call of GetSimilarSeries
func (mr *MockRepositoryMockRecorder) GetSimilarSeries(sid interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetSimilarSeries", reflect.TypeOf((*MockRepository)(nil).GetSimilarSeries), sid)
}
