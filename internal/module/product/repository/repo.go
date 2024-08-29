package repository

import (
	"codebase-app/internal/module/product/entity"
	"codebase-app/internal/module/product/ports"
	"context"

	"github.com/jmoiron/sqlx"
	"github.com/rs/zerolog/log"
)

var _ ports.ProductRepository = &productRepository{}

type productRepository struct {
	db *sqlx.DB
}

func NewProductRepository(db *sqlx.DB) *productRepository {
	return &productRepository{
		db: db,
	}
}

func (r *productRepository) GetProducts(ctx context.Context, req *entity.ProductsRequest) (*entity.ProductsResponse, error) {
	type dao struct {
		TotalData int `db:"total_data"`
		entity.ProductItem
	}

	var (
		resp = new(entity.ProductsResponse)
		data = make([]dao, 0, req.Paginate)
	)
	resp.Items = make([]entity.ProductItem, 0, req.Paginate)

	query := `
		SELECT
			COUNT(p.id) OVER() as total_data,
			p.id,
			p.name,
			p.description,
			p.price,
			p.stock,
			p.user_id,
			s.name AS "shop.name",
			s.description AS "shop.description",
			s.terms AS "shop.terms",
			c.name AS "category.name"
		FROM products p
		INNER JOIN shops s ON p.shop_id = s.id
		INNER JOIN categories c ON p.category_id = c.id
		WHERE
			p.user_id = ?
			AND p.deleted_at IS NULL
	`

	args := []interface{}{req.UserId}

	if req.CategoryId != "" {
		query += " AND c.id = ?"
		args = append(args, req.CategoryId)
	}

	if req.MinPrice != nil {
		query += " AND p.price >= ?"
		args = append(args, *req.MinPrice)
	}

	if req.MaxPrice != nil {
		query += " AND p.price <= ?"
		args = append(args, *req.MaxPrice)
	}

	if req.SearchQuery != "" {
		query += " AND (p.name LIKE ? OR p.description LIKE ?)"
		searchTerm := "%" + req.SearchQuery + "%"
		args = append(args, searchTerm, searchTerm)
	}

	if req.IsAvailable {
		query += " AND p.stock > 0"
	}

	query += "ORDER BY p.created_at DESC LIMIT ? OFFSET ?"
	args = append(args, req.Paginate, req.Paginate*(req.Page-1))

	err := r.db.SelectContext(ctx, &data, r.db.Rebind(query), args...)
	if err != nil {
		log.Error().Err(err).Any("payload", req).Msg("repository::GetProducts - Failed to get products")
		return nil, err
	}

	if len(data) > 0 {
		resp.Meta.TotalData = data[0].TotalData
	}

	for _, d := range data {
		resp.Items = append(resp.Items, d.ProductItem)
	}

	resp.Meta.CountTotalPage(req.Page, req.Paginate, resp.Meta.TotalData)

	return resp, nil
}

func (r *productRepository) CreateProduct(ctx context.Context, req *entity.CreateProductRequest) (*entity.CreateProductResponse, error) {
	var resp = new(entity.CreateProductResponse)

	query := `
		INSERT INTO products (shop_id, category_id, name, description, price, stock, user_id)
		VALUES (?, ?, ?, ?, ?, ?, ?) RETURNING id
	`

	err := r.db.QueryRowxContext(ctx, r.db.Rebind(query),
		req.ShopId,
		req.CategoryId,
		req.Name,
		req.Description,
		req.Price,
		req.Stock,
		req.UserId).Scan(&resp.Id)
	if err != nil {
		log.Error().Err(err).Any("payload", req).Msg("repository::CreateProduct - Failed to create product")
		return nil, err
	}

	return resp, nil
}

func (r *productRepository) GetProduct(ctx context.Context, req *entity.GetProductRequest) (*entity.GetProductResponse, error) {
	var resp = new(entity.GetProductResponse)

	query := `
	  SELECT
			p.id, p.name, p.description, p.price, p.stock, p.user_id, s.name AS "shop.name", s.description AS "shop.description", s.terms AS "shop.terms", c.name AS "category.name"
		FROM products p
		INNER JOIN shops s ON p.shop_id = s.id
		INNER JOIN categories c ON p.category_id = c.id
		WHERE p.id = ? AND p.deleted_at IS NULL
	`

	err := r.db.QueryRowxContext(ctx, r.db.Rebind(query), req.Id).StructScan(resp)
	if err != nil {
		log.Error().Err(err).Any("payload", req).Msg("repository::GetProduct - Failed to get product")
		return nil, err
	}

	return resp, nil
}

func (r *productRepository) UpdateProduct(ctx context.Context, req *entity.UpdateProductRequest) (*entity.UpdateProductResponse, error) {
	var resp = new(entity.UpdateProductResponse)

	query := `
		UPDATE products
		SET shop_id = ?, category_id = ?, name = ?, description = ?, price = ?, stock = ?, updated_at = NOW()
		WHERE id = ? AND user_id = ?
		RETURNING id
	`

	err := r.db.QueryRowxContext(ctx, r.db.Rebind(query),
		req.ShopId,
		req.CategoryId,
		req.Name,
		req.Description,
		req.Price,
		req.Stock,
		req.Id,
		req.UserId).Scan(&resp.Id)
	if err != nil {
		log.Error().Err(err).Any("payload", req).Msg("repository::UpdateProduct - Failed to update product")
		return nil, err
	}

	return resp, nil
}

func (r *productRepository) DeleteProduct(ctx context.Context, req *entity.DeleteProductRequest) error {
	query := `
		UPDATE products
		SET deleted_at = NOW()
		WHERE id = ? AND user_id = ? AND deleted_at IS NULL
	`

	_, err := r.db.ExecContext(ctx, r.db.Rebind(query), req.Id, req.UserId)
	if err != nil {
		log.Error().Err(err).Any("payload", req).Msg("repository::DeleteProduct - Failed to delete product")
		return err
	}

	return nil
}
