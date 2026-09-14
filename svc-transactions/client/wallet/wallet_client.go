package wallet

import (
	"context"

	"svc-transactions/proto"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type Client struct {
	conn   *grpc.ClientConn
	client proto.WalletServiceClient
}

func NewClient(addr string) (*Client, error) {
	conn, err := grpc.NewClient(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}
	return &Client{
		conn:   conn,
		client: proto.NewWalletServiceClient(conn),
	}, nil
}

// ApplyBalanceChange calls svc-wallet over gRPC.
// requestID is the idempotency key — svc-wallet uses it to recognise a retry
// and replay the original result instead of moving money twice.
func (c *Client) ApplyBalanceChange(ctx context.Context, walletID, changeType string, amount int64, requestID string) (*proto.BalanceChangeResponse, error) {
	return c.client.ApplyBalanceChange(ctx, &proto.BalanceChangeRequest{
		WalletId:  walletID,
		Type:      changeType,
		Amount:    amount,
		RequestId: requestID,
	})
}
