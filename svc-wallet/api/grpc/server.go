package grpc

import (
	"context"
	"svc-wallet/internal/wallet"
	"svc-wallet/proto"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
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

func (s *GRPCServer) ApplyBalanceChange(ctx context.Context, req *proto.BalanceChangeRequest) (*proto.BalanceChangeResponse, error) {
	updatedWallet ,balanceBefore , err := s.service.ApplyBalanceChange(ctx,req.WalletId,req.Type,req.Amount)
	if err != nil {
		return nil , status.Error(codes.Internal ,err.Error())
	}

	return &proto.BalanceChangeResponse{
		WalletId: updatedWallet.ID.Hex(),
		Type: req.Type,
		Amount: req.Amount,
		BalanceBefore: balanceBefore,
		BalanceAfter: updatedWallet.Balance,
	}, nil
}