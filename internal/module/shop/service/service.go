package service

import (
	"codebase-app/internal/module/shop/entity"
	"codebase-app/internal/module/shop/ports"
	"context"
)

var _ ports.ShopService = &shopService{}

type shopService struct {
	repo ports.ShopRepository
}

func NewShopService(repo ports.ShopRepository) *shopService {
	return &shopService{
		repo: repo,
	}
}

func (s *shopService) CreateShop(ctx context.Context, req *entity.CreateShopRequest) (*entity.CreateShopResponse, error) {
	return s.repo.CreateShop(ctx, req)
}

func (s *shopService) GetShop(ctx context.Context, req *entity.GetShopRequest) (*entity.GetShopResponse, error) {
	return s.repo.GetShop(ctx, req)
}

func (s *shopService) DeleteShop(ctx context.Context, req *entity.DeleteShopRequest) error {
	return s.repo.DeleteShop(ctx, req)
}

func (s *shopService) UpdateShop(ctx context.Context, req *entity.UpdateShopRequest) (*entity.UpdateShopResponse, error) {
	return s.repo.UpdateShop(ctx, req)
}

func (s *shopService) GetShops(ctx context.Context, req *entity.ShopsRequest) (*entity.ShopsResponse, error) {
	return s.repo.GetShops(ctx, req)
}
