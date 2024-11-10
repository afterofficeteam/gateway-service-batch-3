package cart

import (
	"context"
	"errors"
	"fmt"
	"gateway-service/proto/cart"
)

type svc struct {
	client cart.CartServiceClient
}

func NewSvc(client cart.CartServiceClient) *svc {
	return &svc{client: client}
}

type CartSvc interface {
	InsertCart(ctx context.Context, req *cart.CartInsertRequest) (*cart.CartInsertResponse, error)
}

func (s *svc) InsertCart(ctx context.Context, req *cart.CartInsertRequest) (*cart.CartInsertResponse, error) {
	resp, err := s.client.InsertCart(ctx, req)
	if err != nil {
		return nil, errors.Join(fmt.Errorf("failed to insert cart: %v", err))
	}

	return resp, nil
}
