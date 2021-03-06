// Code generated by MockGen. DO NOT EDIT.
// Source: ./application/review/usecase.go

// Package mock_review is a generated GoMock package.
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

// Add mocks base method
func (m *MockUseCase) Add(review *models.Review) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Add", review)
	ret0, _ := ret[0].(error)
	return ret0
}

// Add indicates an expected call of Add
func (mr *MockUseCaseMockRecorder) Add(review interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Add", reflect.TypeOf((*MockUseCase)(nil).Add), review)
}

// GetByProductId mocks base method
func (m *MockUseCase) GetByProductId(id uint) ([]models.Review, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetByProductId", id)
	ret0, _ := ret[0].([]models.Review)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetByProductId indicates an expected call of GetByProductId
func (mr *MockUseCaseMockRecorder) GetByProductId(id interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetByProductId", reflect.TypeOf((*MockUseCase)(nil).GetByProductId), id)
}

// GetReview mocks base method
func (m *MockUseCase) GetReview(productId, userId uint) (*models.Review, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetReview", productId, userId)
	ret0, _ := ret[0].(*models.Review)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetReview indicates an expected call of GetReview
func (mr *MockUseCaseMockRecorder) GetReview(productId, userId interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetReview", reflect.TypeOf((*MockUseCase)(nil).GetReview), productId, userId)
}
