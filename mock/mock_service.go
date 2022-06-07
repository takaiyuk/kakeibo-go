// Code generated by MockGen. DO NOT EDIT.
// Source: service.go

// Package mock is a generated GoMock package.
package mock

import (
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
	pkg "github.com/takaiyuk/kakeibo-go/pkg"
)

// MockInterfaceService is a mock of InterfaceService interface.
type MockInterfaceService struct {
	ctrl     *gomock.Controller
	recorder *MockInterfaceServiceMockRecorder
}

// MockInterfaceServiceMockRecorder is the mock recorder for MockInterfaceService.
type MockInterfaceServiceMockRecorder struct {
	mock *MockInterfaceService
}

// NewMockInterfaceService creates a new mock instance.
func NewMockInterfaceService(ctrl *gomock.Controller) *MockInterfaceService {
	mock := &MockInterfaceService{ctrl: ctrl}
	mock.recorder = &MockInterfaceServiceMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockInterfaceService) EXPECT() *MockInterfaceServiceMockRecorder {
	return m.recorder
}

// GetSlackMessages mocks base method.
func (m *MockInterfaceService) GetSlackMessages(arg0 *pkg.Config, arg1 *pkg.FilterSlackMessagesOptions) ([]*pkg.SlackMessage, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetSlackMessages", arg0, arg1)
	ret0, _ := ret[0].([]*pkg.SlackMessage)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetSlackMessages indicates an expected call of GetSlackMessages.
func (mr *MockInterfaceServiceMockRecorder) GetSlackMessages(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetSlackMessages", reflect.TypeOf((*MockInterfaceService)(nil).GetSlackMessages), arg0, arg1)
}

// PostIFTTTWebhook mocks base method.
func (m *MockInterfaceService) PostIFTTTWebhook(arg0 *pkg.Config, arg1 []*pkg.SlackMessage) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "PostIFTTTWebhook", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// PostIFTTTWebhook indicates an expected call of PostIFTTTWebhook.
func (mr *MockInterfaceServiceMockRecorder) PostIFTTTWebhook(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "PostIFTTTWebhook", reflect.TypeOf((*MockInterfaceService)(nil).PostIFTTTWebhook), arg0, arg1)
}