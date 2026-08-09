package grpc

import (
	"svc-user/internal/user"
	"svc-user/proto"
)

type GRPCServer struct {
	proto.UnimplementedUserServiceServer
	service *user.Service
}

func NewGRPCServer(service *user.Service) *GRPCServer {
	return &GRPCServer{
		service: service,
	}
}