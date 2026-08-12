package grpc


import (
	"context"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"svc-wallet/internal/wallet"
	"svc-wallet/proto"	
)

type GRPCServer struct {
	proto.UnimplementedWalletServiceServer
	service *wallet.Service
}

func NewGRPCServer(service *wallet.Service) *GRPCServer{
	return &GRPCServer{
		service: service,
	}
}

