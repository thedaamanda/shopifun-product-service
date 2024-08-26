package service

import (
	integOauth "codebase-app/internal/integration/oauth2google"
	oauthgoogleent "codebase-app/internal/integration/oauth2google/entity"
	"codebase-app/internal/module/user/entity"
	"codebase-app/internal/module/user/ports"
	"codebase-app/pkg"
	"codebase-app/pkg/errmsg"
	"codebase-app/pkg/jwthandler"
	"context"
	"time"

	"github.com/rs/zerolog/log"
)

var _ ports.UserService = &userService{}

type userService struct {
	repo ports.UserRepository
	o    integOauth.Oauth2googleContract
}

func NewUserService(repo ports.UserRepository, o integOauth.Oauth2googleContract) *userService {
	return &userService{
		repo: repo,
		o:    o,
	}
}

func (s *userService) Register(ctx context.Context, req *entity.RegisterRequest) (*entity.RegisterResponse, error) {

	hashed, err := pkg.HashPassword(req.Password)
	if err != nil {
		log.Error().Err(err).Any("payload", req).Msg("service::Register - Failed to hash password")
		return nil, errmsg.NewCustomErrors(500, errmsg.WithMessage("Gagal menghash password"))
	}

	req.HassedPassword = hashed

	result, err := s.repo.Register(ctx, req)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (s *userService) Login(ctx context.Context, req *entity.LoginRequest) (*entity.LoginResponse, error) {
	var res = new(entity.LoginResponse)

	user, err := s.repo.FindByEmail(ctx, req.Email)
	if err != nil {
		return nil, err
	}

	if !pkg.ComparePassword(user.Pass, req.Password) {
		log.Warn().Any("payload", req).Msg("service::Login - Password not match")
		return nil, errmsg.NewCustomErrors(401, errmsg.WithMessage("Email atau password salah"))
	}

	token, err := jwthandler.GenerateTokenString(jwthandler.CostumClaimsPayload{
		UserId:          user.Id,
		Role:            user.Role,
		TokenExpiration: time.Now().Add(time.Hour * 24),
	})
	if err != nil {
		return nil, err
	}

	res.Token = token
	return res, nil
}

func (s *userService) Profile(ctx context.Context, req *entity.ProfileRequest) (*entity.ProfileResponse, error) {
	user, err := s.repo.FindById(ctx, req.UserId)
	if err != nil {
		return nil, err
	}

	return user, nil

}

func (s *userService) GetOauthGoogleUrl(ctx context.Context) (string, error) {
	url := s.o.GetUrl("state")

	return url, nil
}

func (s *userService) LoginGoogle(ctx context.Context, req *oauthgoogleent.UserInfoResponse) (*entity.LoginResponse, error) {
	var res = new(entity.LoginResponse)

	user, err := s.repo.FindByEmail(ctx, req.Email)
	if err != nil {
		if errCostum, ok := err.(*errmsg.CustomError); ok {
			if errCostum.Code != 400 {
				return nil, err
			} else {
				// TODO: logic to register user and save to database
				return nil, errmsg.NewCustomErrors(404, errmsg.WithMessage("Email belum terdaftar"))
			}
		}
		return nil, err
	}

	token, err := jwthandler.GenerateTokenString(jwthandler.CostumClaimsPayload{
		UserId:          user.Id,
		Role:            user.Role,
		TokenExpiration: time.Now().Add(time.Hour * 24),
	})
	if err != nil {
		return nil, err
	}

	res.Token = token
	return res, nil
}
