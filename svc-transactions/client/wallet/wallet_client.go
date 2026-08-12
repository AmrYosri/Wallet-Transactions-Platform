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
	conn , err := grpc.NewClient(addr,grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil{
		return nil , err
	}
	return &Client{
		conn: conn,
		client : proto.NewWalletServiceClient(conn),
	},nil
}

func (c *Client) ApplyBalanceChange(ctx context.Context, walletID, changeType string, amount int64) (*proto.BalanceChangeResponse, error) {
	return c.client.ApplyBalanceChange(ctx , &proto.BalanceChangeRequest{WalletId: walletID, Type: changeType, Amount: amount})

}