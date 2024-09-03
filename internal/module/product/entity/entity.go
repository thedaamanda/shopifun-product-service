package entity

import "codebase-app/pkg/types"

type Shop struct {
	Name        string `json:"name" db:"name"`
	Description string `json:"description" db:"description"`
	Terms       string `json:"terms" db:"terms"`
}

type Category struct {
	Name string `json:"name" db:"name"`
}

type Brand struct {
	Name string `json:"name" db:"name"`
}

type ProductsRequest struct {
	UserId      string   `prop:"user_id" validate:"uuid"`
	Page        int      `query:"page" validate:"required"`
	Paginate    int      `query:"paginate" validate:"required"`
	CategoryId  string   `query:"category_id" validate:"omitempty,uuid"`
	BrandId     string   `query:"brand_id" validate:"omitempty,uuid"`
	MinPrice    *float64 `query:"min_price" validate:"omitempty,numeric,min=0"`
	MaxPrice    *float64 `query:"max_price" validate:"omitempty,numeric,min=0"`
	SearchQuery string   `query:"search_query" validate:"omitempty,min=3,max=255"`
	IsAvailable bool     `query:"is_available" validate:"omitempty"`
}

func (r *ProductsRequest) SetDefault() {
	if r.Page < 1 {
		r.Page = 1
	}

	if r.Paginate < 1 {
		r.Paginate = 10
	}
}

type ProductItem struct {
	Id          string   `json:"id" db:"id"`
	Name        string   `json:"name" db:"name"`
	Description string   `json:"description" db:"description"`
	Price       float64  `json:"price" db:"price"`
	Stock       int      `json:"stock" db:"stock"`
	UserId      string   `json:"user_id" db:"user_id"`
	Category    Category `json:"category"`
	Shop        Shop     `json:"shop"`
	Brand       Brand    `json:"brand"`
}

type ProductsResponse struct {
	Items []ProductItem `json:"items"`
	Meta  types.Meta    `json:"meta"`
}

type CreateProductRequest struct {
	ShopId      string  `json:"shop_id" validate:"required,uuid" db:"shop_id"`
	CategoryId  string  `json:"category_id" validate:"required,uuid" db:"category_id"`
	BrandId     string  `json:"brand_id" validate:"required,uuid" db:"brand_id"`
	Name        string  `json:"name" validate:"required" db:"name"`
	Description string  `json:"description" validate:"required,max=255" db:"description"`
	Price       float64 `json:"price" validate:"required" db:"price"`
	Stock       int     `json:"stock" validate:"required,numeric" db:"stock"`
	UserId      string  `json:"user_id" validate:"uuid" db:"user_id"`
}

type CreateProductResponse struct {
	Id string `json:"id" db:"id"`
}

type GetProductRequest struct {
	Id string `validate:"uuid" db:"id"`
}

type GetProductResponse struct {
	Id          string   `json:"id" db:"id"`
	Name        string   `json:"name" db:"name"`
	Description string   `json:"description" db:"description"`
	Price       float64  `json:"price" db:"price"`
	Stock       int      `json:"stock" db:"stock"`
	UserId      string   `json:"user_id" db:"user_id"`
	Category    Category `json:"category"`
	Shop        Shop     `json:"shop"`
	Brand       Brand    `json:"brand"`
}

type UpdateProductRequest struct {
	Id          string  `params:"id" validate:"uuid" db:"id"`
	ShopId      string  `json:"shop_id" validate:"required,uuid" db:"shop_id"`
	CategoryId  string  `json:"category_id" validate:"required,uuid" db:"category_id"`
	BrandId     string  `json:"brand_id" validate:"required,uuid" db:"brand_id"`
	Name        string  `json:"name" validate:"required" db:"name"`
	Description string  `json:"description" validate:"required,max=255" db:"description"`
	Price       float64 `json:"price" validate:"required" db:"price"`
	Stock       int     `json:"stock" validate:"required,numeric" db:"stock"`
	UserId      string  `json:"user_id" validate:"uuid" db:"user_id"`
}

type UpdateProductResponse struct {
	Id string `json:"id" db:"id"`
}

type DeleteProductRequest struct {
	UserId string `prop:"user_id" validate:"uuid" db:"user_id"`

	Id string `validate:"uuid" db:"id"`
}
