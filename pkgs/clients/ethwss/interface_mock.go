// Code generated by MockGen. DO NOT EDIT.
// Source: interface.go
//
// Generated by this command:
//
//	mockgen -source=interface.go -destination=interface_mock.go -package=ethwss
//

// Package ethwss is a generated GoMock package.
package ethwss

import (
	context "context"
	reflect "reflect"

	ethereum "github.com/ethereum/go-ethereum"
	types "github.com/ethereum/go-ethereum/core/types"
	gomock "go.uber.org/mock/gomock"
)

// MockGethWssClient is a mock of GethWssClient interface.
type MockGethWssClient struct {
	ctrl     *gomock.Controller
	recorder *MockGethWssClientMockRecorder
	isgomock struct{}
}

// MockGethWssClientMockRecorder is the mock recorder for MockGethWssClient.
type MockGethWssClientMockRecorder struct {
	mock *MockGethWssClient
}

// NewMockGethWssClient creates a new mock instance.
func NewMockGethWssClient(ctrl *gomock.Controller) *MockGethWssClient {
	mock := &MockGethWssClient{ctrl: ctrl}
	mock.recorder = &MockGethWssClientMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockGethWssClient) EXPECT() *MockGethWssClientMockRecorder {
	return m.recorder
}

// Close mocks base method.
func (m *MockGethWssClient) Close() {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "Close")
}

// Close indicates an expected call of Close.
func (mr *MockGethWssClientMockRecorder) Close() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Close", reflect.TypeOf((*MockGethWssClient)(nil).Close))
}

// SubscribeFilterLogs mocks base method.
func (m *MockGethWssClient) SubscribeFilterLogs(ctx context.Context, q ethereum.FilterQuery, ch chan<- types.Log) (ethereum.Subscription, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SubscribeFilterLogs", ctx, q, ch)
	ret0, _ := ret[0].(ethereum.Subscription)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// SubscribeFilterLogs indicates an expected call of SubscribeFilterLogs.
func (mr *MockGethWssClientMockRecorder) SubscribeFilterLogs(ctx, q, ch any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SubscribeFilterLogs", reflect.TypeOf((*MockGethWssClient)(nil).SubscribeFilterLogs), ctx, q, ch)
}
