package pricingsvc

import (
	"context"
	"fmt"
	"housekeeper/api/pricing"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type PricingClient struct {
	client pricing.PricingServiceClient
	conn   *grpc.ClientConn
}

func NewPricingClient(address string) (*PricingClient, error) {
	conn, err := grpc.NewClient(address, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, fmt.Errorf("failed to connect to Pricing Service: %w", err)
	}

	client := pricing.NewPricingServiceClient(conn)
	return &PricingClient{client: client, conn: conn}, nil
}

func (c *PricingClient) GetPrice(ctx context.Context, startDate uint64, requiredHour int32) (*pricing.GetPriceResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	req := &pricing.GetPriceRequest{
		StartDate:    startDate,
		RequiredHour: uint32(requiredHour),
	}
	resp, err := c.client.GetPrice(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("failed to get price from Pricing Service: %w", err)
	}

	return resp, nil
}

func (c *PricingClient) Close() {
	if c.conn != nil {
		c.conn.Close()
	}
}
