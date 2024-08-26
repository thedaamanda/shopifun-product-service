package entity

import "codebase-app/pkg/types"

type CreateShopRequest struct {
	UserId string `validate:"uuid" db:"user_id"`

	Name        string `json:"name" validate:"required" db:"name"`
	Description string `json:"description" validate:"required,max=255" db:"description"`
	Terms       string `json:"terms" validate:"required" db:"terms"`
}

type CreateShopResponse struct {
	Id string `json:"id" db:"id"`
}

type GetShopRequest struct {
	Id string `validate:"uuid" db:"id"`
}

type GetShopResponse struct {
	Name        string `json:"name" db:"name"`
	Description string `json:"description" db:"description"`
	Terms       string `json:"terms" db:"terms"`
}

type DeleteShopRequest struct {
	UserId string `prop:"user_id" validate:"uuid" db:"user_id"`

	Id string `validate:"uuid" db:"id"`
}

type UpdateShopRequest struct {
	UserId string `prop:"user_id" validate:"uuid" db:"user_id"`

	Id          string `params:"id" validate:"uuid" db:"id"`
	Name        string `json:"name" validate:"required" db:"name"`
	Description string `json:"description" validate:"required" db:"description"`
	Terms       string `json:"terms" validate:"required" db:"terms"`
}

type UpdateShopResponse struct {
	Id string `json:"id" db:"id"`
}

type ShopsRequest struct {
	UserId   string `prop:"user_id" validate:"uuid"`
	Page     int    `query:"page" validate:"required"`
	Paginate int    `query:"paginate" validate:"required"`
}

func (r *ShopsRequest) SetDefault() {
	if r.Page < 1 {
		r.Page = 1
	}

	if r.Paginate < 1 {
		r.Paginate = 10
	}
}

type ShopItem struct {
	Id   string `json:"id" db:"id"`
	Name string `json:"name" db:"name"`
}

type ShopsResponse struct {
	Items []ShopItem `json:"items"`
	Meta  types.Meta `json:"meta"`
}
