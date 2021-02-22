// Code generated by MockGen. DO NOT EDIT.
// Source: contracts/build/go/go_opentracing_example/grpc_server/todo/v1/todo_service_grpc.pb.go

// Package todoclientmock is a generated GoMock package.
package todoclientmock

import (
	context "context"
	reflect "reflect"

	v1 "github.com/andream16/go-opentracing-example/contracts/build/go/go_opentracing_example/grpc_server/todo/v1"
	gomock "github.com/golang/mock/gomock"
	grpc "google.golang.org/grpc"
)

// MockTodoServiceClient is a mock of TodoServiceClient interface.
type MockTodoServiceClient struct {
	ctrl     *gomock.Controller
	recorder *MockTodoServiceClientMockRecorder
}

// MockTodoServiceClientMockRecorder is the mock recorder for MockTodoServiceClient.
type MockTodoServiceClientMockRecorder struct {
	mock *MockTodoServiceClient
}

// NewMockTodoServiceClient creates a new mock instance.
func NewMockTodoServiceClient(ctrl *gomock.Controller) *MockTodoServiceClient {
	mock := &MockTodoServiceClient{ctrl: ctrl}
	mock.recorder = &MockTodoServiceClientMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockTodoServiceClient) EXPECT() *MockTodoServiceClientMockRecorder {
	return m.recorder
}

// Create mocks base method.
func (m *MockTodoServiceClient) Create(ctx context.Context, in *v1.CreateRequest, opts ...grpc.CallOption) (*v1.CreateResponse, error) {
	m.ctrl.T.Helper()
	varargs := []interface{}{ctx, in}
	for _, a := range opts {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "Create", varargs...)
	ret0, _ := ret[0].(*v1.CreateResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Create indicates an expected call of Create.
func (mr *MockTodoServiceClientMockRecorder) Create(ctx, in interface{}, opts ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{ctx, in}, opts...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Create", reflect.TypeOf((*MockTodoServiceClient)(nil).Create), varargs...)
}

// MockTodoServiceServer is a mock of TodoServiceServer interface.
type MockTodoServiceServer struct {
	ctrl     *gomock.Controller
	recorder *MockTodoServiceServerMockRecorder
}

// MockTodoServiceServerMockRecorder is the mock recorder for MockTodoServiceServer.
type MockTodoServiceServerMockRecorder struct {
	mock *MockTodoServiceServer
}

// NewMockTodoServiceServer creates a new mock instance.
func NewMockTodoServiceServer(ctrl *gomock.Controller) *MockTodoServiceServer {
	mock := &MockTodoServiceServer{ctrl: ctrl}
	mock.recorder = &MockTodoServiceServerMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockTodoServiceServer) EXPECT() *MockTodoServiceServerMockRecorder {
	return m.recorder
}

// Create mocks base method.
func (m *MockTodoServiceServer) Create(arg0 context.Context, arg1 *v1.CreateRequest) (*v1.CreateResponse, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Create", arg0, arg1)
	ret0, _ := ret[0].(*v1.CreateResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Create indicates an expected call of Create.
func (mr *MockTodoServiceServerMockRecorder) Create(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Create", reflect.TypeOf((*MockTodoServiceServer)(nil).Create), arg0, arg1)
}

// MockUnsafeTodoServiceServer is a mock of UnsafeTodoServiceServer interface.
type MockUnsafeTodoServiceServer struct {
	ctrl     *gomock.Controller
	recorder *MockUnsafeTodoServiceServerMockRecorder
}

// MockUnsafeTodoServiceServerMockRecorder is the mock recorder for MockUnsafeTodoServiceServer.
type MockUnsafeTodoServiceServerMockRecorder struct {
	mock *MockUnsafeTodoServiceServer
}

// NewMockUnsafeTodoServiceServer creates a new mock instance.
func NewMockUnsafeTodoServiceServer(ctrl *gomock.Controller) *MockUnsafeTodoServiceServer {
	mock := &MockUnsafeTodoServiceServer{ctrl: ctrl}
	mock.recorder = &MockUnsafeTodoServiceServerMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockUnsafeTodoServiceServer) EXPECT() *MockUnsafeTodoServiceServerMockRecorder {
	return m.recorder
}

// mustEmbedUnimplementedTodoServiceServer mocks base method.
func (m *MockUnsafeTodoServiceServer) mustEmbedUnimplementedTodoServiceServer() {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "mustEmbedUnimplementedTodoServiceServer")
}

// mustEmbedUnimplementedTodoServiceServer indicates an expected call of mustEmbedUnimplementedTodoServiceServer.
func (mr *MockUnsafeTodoServiceServerMockRecorder) mustEmbedUnimplementedTodoServiceServer() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "mustEmbedUnimplementedTodoServiceServer", reflect.TypeOf((*MockUnsafeTodoServiceServer)(nil).mustEmbedUnimplementedTodoServiceServer))
}