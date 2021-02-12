package todo

import (
	"context"
	"fmt"
	"log"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	todov1 "github.com/andream16/go-opentracing-example/contracts/build/go/go_opentracing_example/grpc_server/todo/v1"
)

type Service struct{}

// NewService returns a new Service.
func NewService() Service {
	return Service{}
}

// Creates a new todo.
func (svc Service) Create(ctx context.Context, req *todov1.CreateRequest) (*todov1.CreateResponse, error) {
	if req == nil {
		log.Println("received nil request for creating a todo")
		return nil, status.Error(codes.InvalidArgument, "received nil request for creating a todo")
	}
	log.Println(fmt.Sprintf("creating a new todo with message: %s", req.Message))
	return &todov1.CreateResponse{}, nil
}
