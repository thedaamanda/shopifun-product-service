package rest

import (
	"codebase-app/internal/adapter"
	"codebase-app/internal/infrastructure/config"
	integOauth "codebase-app/internal/integration/oauth2google"
	oauth "codebase-app/internal/integration/oauth2google/entity"
	"codebase-app/internal/middleware"
	"codebase-app/internal/module/user/entity"
	"encoding/json"

	"codebase-app/internal/module/user/ports"
	"codebase-app/internal/module/user/repository"
	"codebase-app/internal/module/user/service"
	"codebase-app/pkg/errmsg"
	"codebase-app/pkg/response"
	"context"
	"net/http"

	"github.com/coreos/go-oidc"

	"github.com/gofiber/fiber/v2"
	"github.com/rs/zerolog/log"
)

type userHandler struct {
	service     ports.UserService
	integration integOauth.Oauth2googleContract
}

func NewUserHandler(o integOauth.Oauth2googleContract) *userHandler {
	var handler = new(userHandler)

	repo := repository.NewUserRepository(adapter.Adapters.ShopeefunPostgres)
	service := service.NewUserService(repo, o)

	handler.integration = o

	handler.service = service

	return handler
}

func (h *userHandler) Register(router fiber.Router) {
	router.Post("/register", h.register)
	router.Post("/login", h.login)
	router.Get("/profile", middleware.AuthBearer, h.profile)
	router.Get("/profile/:user_id", middleware.AuthBearer, h.profileByUserId)

	router.Get("/oauth/google/url", h.oauthGoogleUrl)
	router.Get("/signin/callback", h.callbackSigninGoogle)
}

func (h *userHandler) register(c *fiber.Ctx) error {
	var (
		req = new(entity.RegisterRequest)
		ctx = c.Context()
		v   = adapter.Adapters.Validator
	)

	if err := c.BodyParser(req); err != nil {
		log.Warn().Err(err).Msg("handler::register - Failed to parse request body")
		return c.Status(fiber.StatusBadRequest).JSON(response.Error(err))
	}

	if err := v.Validate(req); err != nil {
		log.Warn().Err(err).Msg("handler::register - Invalid request body")
		code, errs := errmsg.Errors(err, req)
		return c.Status(code).JSON(response.Error(errs))
	}

	res, err := h.service.Register(ctx, req)
	if err != nil {
		code, errs := errmsg.Errors[error](err)
		return c.Status(code).JSON(response.Error(errs))
	}

	return c.Status(fiber.StatusCreated).JSON(response.Success(res, ""))
}

func (h *userHandler) login(c *fiber.Ctx) error {
	var (
		req = new(entity.LoginRequest)
		ctx = c.Context()
		v   = adapter.Adapters.Validator
	)

	if err := c.BodyParser(req); err != nil {
		log.Warn().Err(err).Msg("handler::login - Failed to parse request body")
		return c.Status(fiber.StatusBadRequest).JSON(response.Error(err))
	}

	if err := v.Validate(req); err != nil {
		log.Warn().Err(err).Msg("handler::login - Invalid request body")
		code, errs := errmsg.Errors(err, req)
		return c.Status(code).JSON(response.Error(errs))
	}

	res, err := h.service.Login(ctx, req)
	if err != nil {
		code, errs := errmsg.Errors[error](err)
		return c.Status(code).JSON(response.Error(errs))
	}

	return c.Status(fiber.StatusOK).JSON(response.Success(res, ""))
}

func (h *userHandler) profileByUserId(c *fiber.Ctx) error {
	var (
		req = new(entity.ProfileRequest)
		ctx = c.Context()
		v   = adapter.Adapters.Validator
	)

	req.UserId = c.Params("user_id")

	if err := v.Validate(req); err != nil {
		log.Warn().Err(err).Msg("handler::profileByUserId - Invalid Request")
		code, errs := errmsg.Errors(err, req)
		return c.Status(code).JSON(response.Error(errs))
	}

	res, err := h.service.Profile(ctx, req)
	if err != nil {
		code, errs := errmsg.Errors[error](err)
		return c.Status(code).JSON(response.Error(errs))
	}

	return c.Status(fiber.StatusOK).JSON(response.Success(res, ""))
}

func (h *userHandler) profile(c *fiber.Ctx) error {
	var (
		req = new(entity.ProfileRequest)
		ctx = c.Context()
		l   = middleware.GetLocals(c)
	)

	req.UserId = l.GetUserId()

	res, err := h.service.Profile(ctx, req)
	if err != nil {
		code, errs := errmsg.Errors[error](err)
		return c.Status(code).JSON(response.Error(errs))
	}

	return c.Status(fiber.StatusOK).JSON(response.Success(res, ""))
}

func (h *userHandler) oauthGoogleUrl(c *fiber.Ctx) error {
	return c.Redirect(h.integration.GetUrl("/"), http.StatusTemporaryRedirect)
}

func (h *userHandler) callbackSigninGoogle(c *fiber.Ctx) error {
	var (
		ctx = c.Context()
	)

	state, code := c.FormValue("state"), c.FormValue("code")
	if state == "" && code == "" {
		return c.Status(fiber.StatusBadRequest).JSON(response.Error(errmsg.NewCustomErrors(400, errmsg.WithMessage("Invalid request"))))
	}

	token, err := h.integration.Exchange(ctx, code)
	if err != nil {
		code, errs := errmsg.Errors[error](err)
		return c.Status(code).JSON(response.Error(errs))
	}

	provider, err := oidc.NewProvider(ctx, "https://accounts.google.com")
	if err != nil {
		code, errs := errmsg.Errors[error](err)
		return c.Status(code).JSON(response.Error(errs))
	}

	verifier := provider.Verifier(&oidc.Config{
		ClientID: config.Envs.Oauth.Google.ClientId,
	})
	_, err = verifier.Verify(context.Background(), token.Extra("id_token").(string))
	if err != nil {
		code, errs := errmsg.Errors[error](err)
		return c.Status(code).JSON(response.Error(errs))
	}

	result, err := http.Get("https://www.googleapis.com/oauth2/v2/userinfo?access_token=" + token.AccessToken)
	if err != nil {
		code, errs := errmsg.Errors[error](err)
		return c.Status(code).JSON(response.Error(errs))
	}
	defer result.Body.Close()

	var userInfo oauth.UserInfoResponse
	if err := json.NewDecoder(result.Body).Decode(&userInfo); err != nil {
		code, errs := errmsg.Errors[error](err)
		return c.Status(code).JSON(response.Error(errs))
	}

	res, err := h.service.LoginGoogle(ctx, &userInfo)
	if err != nil {
		code, errs := errmsg.Errors[error](err)
		return c.Status(code).JSON(response.Error(errs))
	}

	return c.Status(fiber.StatusOK).JSON(response.Success(res, ""))
}

// Convert dari PRD ke user story
// Jelaskan diagram dan alur based on user story
// Jelaskan based on diagram
