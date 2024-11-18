package cart

import (
	"context"
	"errors"
	"fmt"
	"gateway-service/proto/cart"
)

type svc struct {
	server cart.CartServiceClient
}

func NewSvc(client cart.CartServiceClient) *svc {
	return &svc{server: client}
}

type CartSvc interface {
	InsertCart(ctx context.Context, req *cart.CartInsertRequest) (*cart.CartInsertResponse, error)
	GetDetail(ctx context.Context, req *cart.CartDetailRequest) (*cart.CartDetailResponse, error)
}

func (s *svc) InsertCart(ctx context.Context, req *cart.CartInsertRequest) (*cart.CartInsertResponse, error) {
	res, err := s.server.InsertCart(ctx, req)
	if err != nil {
		return nil, errors.Join(fmt.Errorf("failed to insert cart: %v", err))
	}

	return res, nil
}

func (s *svc) GetDetail(ctx context.Context, req *cart.CartDetailRequest) (*cart.CartDetailResponse, error) {
	res, err := s.server.DetailCart(ctx, req)
	if err != nil {
		return nil, errors.Join(fmt.Errorf("failed to get detail cart: %v", err))
	}

	return res, nil
}
