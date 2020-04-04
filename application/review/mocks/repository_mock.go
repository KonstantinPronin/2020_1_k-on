// Code generated by MockGen. DO NOT EDIT.
// Source: ./application/review/repository.go

// Package mock_review is a generated GoMock package.
package mocks

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

// Add mocks base method
func (m *MockRepository) Add(review *models.Review) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Add", review)
	ret0, _ := ret[0].(error)
	return ret0
}

// Add indicates an expected call of Add
func (mr *MockRepositoryMockRecorder) Add(review interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Add", reflect.TypeOf((*MockRepository)(nil).Add), review)
}

// GetByProductId mocks base method
func (m *MockRepository) GetByProductId(id uint) ([]models.Review, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetByProductId", id)
	ret0, _ := ret[0].([]models.Review)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetByProductId indicates an expected call of GetByProductId
func (mr *MockRepositoryMockRecorder) GetByProductId(id interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetByProductId", reflect.TypeOf((*MockRepository)(nil).GetByProductId), id)
}