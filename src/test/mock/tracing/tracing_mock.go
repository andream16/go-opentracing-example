// Code generated by MockGen. DO NOT EDIT.
// Source: src/shared/tracing/tracing.go

// Package tracingmock is a generated GoMock package.
package tracingmock

import (
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
	opentracing "github.com/opentracing/opentracing-go"
)

// MockTracer is a mock of Tracer interface.
type MockTracer struct {
	ctrl     *gomock.Controller
	recorder *MockTracerMockRecorder
}

// MockTracerMockRecorder is the mock recorder for MockTracer.
type MockTracerMockRecorder struct {
	mock *MockTracer
}

// NewMockTracer creates a new mock instance.
func NewMockTracer(ctrl *gomock.Controller) *MockTracer {
	mock := &MockTracer{ctrl: ctrl}
	mock.recorder = &MockTracerMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockTracer) EXPECT() *MockTracerMockRecorder {
	return m.recorder
}

// Close mocks base method.
func (m *MockTracer) Close() error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Close")
	ret0, _ := ret[0].(error)
	return ret0
}

// Close indicates an expected call of Close.
func (mr *MockTracerMockRecorder) Close() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Close", reflect.TypeOf((*MockTracer)(nil).Close))
}

// Extract mocks base method.
func (m *MockTracer) Extract(format, carrier interface{}) (opentracing.SpanContext, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Extract", format, carrier)
	ret0, _ := ret[0].(opentracing.SpanContext)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Extract indicates an expected call of Extract.
func (mr *MockTracerMockRecorder) Extract(format, carrier interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Extract", reflect.TypeOf((*MockTracer)(nil).Extract), format, carrier)
}

// Inject mocks base method.
func (m *MockTracer) Inject(sm opentracing.SpanContext, format, carrier interface{}) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Inject", sm, format, carrier)
	ret0, _ := ret[0].(error)
	return ret0
}

// Inject indicates an expected call of Inject.
func (mr *MockTracerMockRecorder) Inject(sm, format, carrier interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Inject", reflect.TypeOf((*MockTracer)(nil).Inject), sm, format, carrier)
}

// StartSpan mocks base method.
func (m *MockTracer) StartSpan(operationName string, opts ...opentracing.StartSpanOption) opentracing.Span {
	m.ctrl.T.Helper()
	varargs := []interface{}{operationName}
	for _, a := range opts {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "StartSpan", varargs...)
	ret0, _ := ret[0].(opentracing.Span)
	return ret0
}

// StartSpan indicates an expected call of StartSpan.
func (mr *MockTracerMockRecorder) StartSpan(operationName interface{}, opts ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{operationName}, opts...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "StartSpan", reflect.TypeOf((*MockTracer)(nil).StartSpan), varargs...)
}
