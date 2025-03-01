package sendsvc

import (
	"context"
	"fmt"
	"housekeeper/api/send"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type SendClient struct {
	client send.SendServiceClient
	conn   *grpc.ClientConn
}

func NewSendClient(address string) (*SendClient, error) {
	conn, err := grpc.NewClient(address, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, fmt.Errorf("failed to connect to Send Service: %w", err)
	}

	client := send.NewSendServiceClient(conn)
	return &SendClient{client: client, conn: conn}, nil
}

func (c *SendClient) SendJob(ctx context.Context, req *send.SendJobRequest) error {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	_, err := c.client.SendJob(ctx, req)
	if err != nil {
		return fmt.Errorf("failed to send job to Send Service: %w", err)
	}

	return nil
}

func (c *SendClient) Close() {
	if c.conn != nil {
		c.conn.Close()
	}
}
