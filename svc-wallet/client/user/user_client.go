package user

import (
	"context"

	"svc-wallet/proto"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type Client struct {
	conn   *grpc.ClientConn
	client proto.UserServiceClient
}

func NewClient(addr string) (*Client, error) {
	conn, err := grpc.NewClient(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}
	return &Client{
		conn:   conn,
		client: proto.NewUserServiceClient(conn),
	}, nil
}

func (c *Client) GetUserByNationalID(ctx context.Context, nationalID string) (*proto.UserResponse, error) {
	return c.client.GetUserByNationalID(ctx, &proto.GetUserRequest{NationalId: nationalID})
}