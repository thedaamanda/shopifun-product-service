package service

import (
	"codebase-app/internal/module/product/entity"
	"codebase-app/internal/module/product/ports"
	"codebase-app/pkg/errmsg"
	"context"

	"github.com/rs/zerolog/log"
)

var _ ports.ProductService = &productService{}

type productService struct {
	repo ports.ProductRepository
}

func NewProductService(repo ports.ProductRepository) *productService {
	return &productService{
		repo: repo,
	}
}

func (s *productService) GetProducts(ctx context.Context, req *entity.ProductsRequest) (*entity.ProductsResponse, error) {
	res, err := s.repo.GetProducts(ctx, req)
	if err != nil {
		return res, err
	}

	if len(res.Items) == 0 {
		log.Warn().Any("payload", req).Msg("service: Products not found")
		return res, errmsg.NewCustomErrors(404, errmsg.WithMessage("Products not found"))
	}

	return res, nil
}

func (s *productService) CreateProduct(ctx context.Context, req *entity.CreateProductRequest) (*entity.CreateProductResponse, error) {
	var res *entity.CreateProductResponse

	isShopOwner, err := s.repo.IsShopOwner(ctx, req.UserId, req.ShopId)
	if err != nil {
		return res, err
	}

	if !isShopOwner {
		log.Warn().Any("payload", req).Msg("service: User is not shop owner")
		return res, errmsg.NewCustomErrors(403, errmsg.WithMessage("User is not shop owner"))
	}

	res, err = s.repo.CreateProduct(ctx, req)
	if err != nil {
		return res, err
	}

	return res, nil
}

func (s *productService) GetProduct(ctx context.Context, req *entity.GetProductRequest) (*entity.GetProductResponse, error) {
	return s.repo.GetProduct(ctx, req)
}

func (s *productService) UpdateProduct(ctx context.Context, req *entity.UpdateProductRequest) (*entity.UpdateProductResponse, error) {
	var res *entity.UpdateProductResponse

	isProductOwner, err := s.repo.IsProductOwner(ctx, req.UserId, req.Id)
	if err != nil {
		return res, err
	}

	if !isProductOwner {
		log.Warn().Any("payload", req).Msg("service: User is not product owner")
		return res, errmsg.NewCustomErrors(403, errmsg.WithMessage("User is not product owner"))
	}

	res, err = s.repo.UpdateProduct(ctx, req)
	if err != nil {
		return res, err
	}

	return res, nil
}

func (s *productService) DeleteProduct(ctx context.Context, req *entity.DeleteProductRequest) error {
	isProductOwner, err := s.repo.IsProductOwner(ctx, req.UserId, req.Id)
	if err != nil {
		return err
	}

	if !isProductOwner {
		log.Warn().Any("payload", req).Msg("service: User is not product owner")
		return errmsg.NewCustomErrors(403, errmsg.WithMessage("User is not product owner"))
	}

	return s.repo.DeleteProduct(ctx, req)
}
