package grpc

import (
	"context"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

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

func (s *GRPCServer) GetUserByNationalID(ctx context.Context, req *proto.GetUserRequest) (*proto.UserResponse, error) {
	foundUser, err := s.service.GetUserByNationalID(ctx, req.NationalId)
	if err != nil {
		return nil, status.Error(codes.NotFound, "user not found")
	}

	return &proto.UserResponse{
		Id:           foundUser.ID.Hex(),
		FirstName:    foundUser.FirstName,
		LastName:     foundUser.LastName,
		NationalId:   foundUser.NationalID,
		PhoneNumbers: foundUser.PhoneNumbers,
		Status:       foundUser.Status,
	}, nil
}

func (s *GRPCServer) AddPhoneNumber(ctx context.Context, req *proto.AddPhoneNumberRequest) (*proto.AddPhoneNumberResponse, error) {
	err := s.service.AddPhoneNumber(ctx, req.NationalId, req.Phone)
	if err != nil {
		return nil, status.Error(codes.Internal, "failed to add phone number")
	}

	return &proto.AddPhoneNumberResponse{
		Success: true,
	}, nil
}
