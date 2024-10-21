package service

import (
	"errors"
	"testing"

	"codebase-app/internal/module/shop/entity"
	"codebase-app/internal/module/shop/ports"
	mockPort "codebase-app/mock/module/shop/ports"
	"codebase-app/pkg/types"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	"golang.org/x/net/context"
)

type MockService struct {
	mock.Mock
}

type ServiceList struct {
	suite.Suite
	mockShopRepo *mockPort.MockShopRepo
	service      ports.ShopService

	mockCreateShopReq       *entity.CreateShopRequest
	mockUpdateShopReq       *entity.UpdateShopRequest
	mockGetShopsReq         *entity.ShopsRequest
	mockGetShopRes          entity.ShopsResponse
	mockGetShopEmptyShopRes entity.ShopsResponse
}

func (suite *ServiceList) SetupTest() {
	suite.mockShopRepo = new(mockPort.MockShopRepo)
	suite.service = NewShopService(suite.mockShopRepo)
	suite.mockCreateShopReq = &entity.CreateShopRequest{
		UserId: "1",
		Name:   "Shop 1",
	}
	suite.mockUpdateShopReq = &entity.UpdateShopRequest{
		UserId: "1",
		Id:     "2",
		Name:   "Shop 1",
	}
	suite.mockGetShopsReq = &entity.ShopsRequest{
		UserId: "1",
	}
	suite.mockGetShopRes = entity.ShopsResponse{
		Items: []entity.ShopItem{
			{
				Id:   "1",
				Name: "Shop 1",
			},
		},
		Meta: types.Meta{
			TotalData: 1,
			TotalPage: 1,
			Page:      1,
			Paginate:  10,
		},
	}
	suite.mockGetShopEmptyShopRes = entity.ShopsResponse{
		Items: []entity.ShopItem{},
		Meta: types.Meta{
			TotalData: 0,
			TotalPage: 1,
			Page:      1,
			Paginate:  10,
		},
	}
}

func (suite *ServiceList) TestCreateShop_Success() {
	ctx := context.Background()
	req := suite.mockCreateShopReq
	suite.mockShopRepo.On("CreateShop", ctx, req).Return(entity.CreateShopResponse{}, nil)
	_, err := suite.service.CreateShop(ctx, req)

	suite.Equal(nil, err)
}

func (suite *ServiceList) TestCreateShop_Failed() {
	ctx := context.Background()
	req := suite.mockCreateShopReq
	suite.mockShopRepo.On("CreateShop", ctx, req).Return(mock.Anything, errors.New(mock.Anything))
	_, err := suite.service.CreateShop(ctx, req)

	suite.Equal(errors.New(mock.Anything), err)
}

func (suite *ServiceList) TestGetShop_Success() {
	ctx := context.Background()
	req := &entity.GetShopRequest{
		Id: "1",
	}
	suite.mockShopRepo.On("GetShop", ctx, req).Return(&entity.GetShopResponse{}, nil)
	_, err := suite.service.GetShop(ctx, req)

	suite.Equal(nil, err)
}

func (suite *ServiceList) TestGetShop_GetShopRepoError() {
	ctx := context.Background()
	req := &entity.GetShopRequest{
		Id: "1",
	}
	suite.mockShopRepo.On("GetShop", ctx, req).Return(mock.Anything, errors.New(mock.Anything))
	_, err := suite.service.GetShop(ctx, req)

	suite.Equal(errors.New(mock.Anything), err)
}

func (suite *ServiceList) TestDeleteShop_Success() {
	ctx := context.Background()
	req := &entity.DeleteShopRequest{
		Id: "1",
	}
	suite.mockShopRepo.On("DeleteShop", ctx, req).Return(nil)
	err := suite.service.DeleteShop(ctx, req)

	suite.Equal(nil, err)
}

func (suite *ServiceList) TestDeleteShop_Failed() {
	ctx := context.Background()
	req := &entity.DeleteShopRequest{
		Id: "1",
	}
	suite.mockShopRepo.On("DeleteShop", ctx, req).Return(errors.New(mock.Anything))
	err := suite.service.DeleteShop(ctx, req)

	suite.Equal(errors.New(mock.Anything), err)
}

func (suite *ServiceList) TestUpdateShop_Success() {
	ctx := context.Background()
	req := suite.mockUpdateShopReq
	suite.mockShopRepo.On("UpdateShop", ctx, req).Return(entity.UpdateShopResponse{}, nil)
	_, err := suite.service.UpdateShop(ctx, req)

	suite.Equal(nil, err)
}

func (suite *ServiceList) TestUpdateShop_Failed() {
	ctx := context.Background()
	req := suite.mockUpdateShopReq
	suite.mockShopRepo.On("UpdateShop", ctx, req).Return(mock.Anything, errors.New(mock.Anything))
	_, err := suite.service.UpdateShop(ctx, req)

	suite.Equal(errors.New(mock.Anything), err)
}

func (suite *ServiceList) TestGetShops_Success() {
	ctx := context.Background()
	req := suite.mockGetShopsReq
	suite.mockShopRepo.On("GetShops", ctx, req).Return(&suite.mockGetShopRes, nil)
	_, err := suite.service.GetShops(ctx, req)

	suite.Equal(nil, err)
}

func (suite *ServiceList) TestGetShops_EmptyShop() {
	ctx := context.Background()
	req := suite.mockGetShopsReq
	suite.mockShopRepo.On("GetShops", ctx, req).Return(&suite.mockGetShopEmptyShopRes, nil)
	_, err := suite.service.GetShops(ctx, req)

	suite.Equal(nil, err)
}

func TestService(t *testing.T) {
	suite.Run(t, new(ServiceList))
}
