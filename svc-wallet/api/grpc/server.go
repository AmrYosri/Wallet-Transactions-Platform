package grpc

import (
	"context"
	"net/http"

	"svc-wallet/internal/wallet"
	"svc-wallet/proto"
	"svc-wallet/util/apperror"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type GRPCServer struct {
	proto.UnimplementedWalletServiceServer
	service *wallet.Service
}

func NewGRPCServer(service *wallet.Service) *GRPCServer {
	return &GRPCServer{
		service: service,
	}
}

func (s *GRPCServer) ApplyBalanceChange(ctx context.Context, req *proto.BalanceChangeRequest) (*proto.BalanceChangeResponse, error) {
	updatedWallet, balanceBefore, err := s.service.ApplyBalanceChange(ctx, req.WalletId, req.Type, req.Amount, req.RequestId)
	if err != nil {
		return nil, toGRPCError(err)
	}

	return &proto.BalanceChangeResponse{
		WalletId:      updatedWallet.ID.Hex(),
		Type:          req.Type,
		Amount:        req.Amount,
		BalanceBefore: balanceBefore,
		BalanceAfter:  updatedWallet.Balance,
	}, nil
}

// toGRPCError maps an AppError onto a gRPC status.
// The domain code travels in the status message so svc-transactions can
// recover it — without that, every failure looks identical on the other side
// and the caller can't tell "insufficient funds" from "no such wallet".
func toGRPCError(err error) error {
	appErr, ok := err.(*apperror.AppError)
	if !ok {
		return status.Error(codes.Internal, apperror.Internal.Code)
	}

	var code codes.Code
	switch appErr.Status {
	case http.StatusNotFound:
		code = codes.NotFound
	case http.StatusConflict:
		code = codes.AlreadyExists
	case http.StatusUnprocessableEntity:
		code = codes.FailedPrecondition
	case http.StatusBadRequest:
		code = codes.InvalidArgument
	default:
		code = codes.Internal
	}

	return status.Error(code, appErr.Code)
}
