// Code generated by MockGen. DO NOT EDIT.
// Source: ./application/person/usecase.go

// Package mock_p is a generated GoMock package.
package mocks

import (
	models "github.com/go-park-mail-ru/2020_1_k-on/application/models"
	gomock "github.com/golang/mock/gomock"
	reflect "reflect"
)

// MockUseCase is a mock of UseCase interface
type MockUseCase struct {
	ctrl     *gomock.Controller
	recorder *MockUseCaseMockRecorder
}

// MockUseCaseMockRecorder is the mock recorder for MockUseCase
type MockUseCaseMockRecorder struct {
	mock *MockUseCase
}

// NewMockUseCase creates a new mock instance
func NewMockUseCase(ctrl *gomock.Controller) *MockUseCase {
	mock := &MockUseCase{ctrl: ctrl}
	mock.recorder = &MockUseCaseMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockUseCase) EXPECT() *MockUseCaseMockRecorder {
	return m.recorder
}

// GetById mocks base method
func (m *MockUseCase) GetById(id uint) (*models.Person, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetById", id)
	ret0, _ := ret[0].(*models.Person)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetById indicates an expected call of GetById
func (mr *MockUseCaseMockRecorder) GetById(id interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetById", reflect.TypeOf((*MockUseCase)(nil).GetById), id)
}

// Add mocks base method
func (m *MockUseCase) Add(p *models.Person) (*models.Person, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Add", p)
	ret0, _ := ret[0].(*models.Person)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Add indicates an expected call of Add
func (mr *MockUseCaseMockRecorder) Add(p interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Add", reflect.TypeOf((*MockUseCase)(nil).Add), p)
}

// Update mocks base method
func (m *MockUseCase) Update(p *models.Person) (*models.Person, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Update", p)
	ret0, _ := ret[0].(*models.Person)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Update indicates an expected call of Update
func (mr *MockUseCaseMockRecorder) Update(p interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Update", reflect.TypeOf((*MockUseCase)(nil).Update), p)
}

// GetActorsForFilm mocks base method
func (m *MockUseCase) GetActorsForFilm(filmId uint) (models.ListPersonArr, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetActorsForFilm", filmId)
	ret0, _ := ret[0].(models.ListPersonArr)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetActorsForFilm indicates an expected call of GetActorsForFilm
func (mr *MockUseCaseMockRecorder) GetActorsForFilm(filmId interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetActorsForFilm", reflect.TypeOf((*MockUseCase)(nil).GetActorsForFilm), filmId)
}

// GetActorsForSeries mocks base method
func (m *MockUseCase) GetActorsForSeries(seriesId uint) (models.ListPersonArr, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetActorsForSeries", seriesId)
	ret0, _ := ret[0].(models.ListPersonArr)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetActorsForSeries indicates an expected call of GetActorsForSeries
func (mr *MockUseCaseMockRecorder) GetActorsForSeries(seriesId interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetActorsForSeries", reflect.TypeOf((*MockUseCase)(nil).GetActorsForSeries), seriesId)
}