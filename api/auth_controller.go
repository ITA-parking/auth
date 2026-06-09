package api

import (
	"auth-service/dto"
	"auth-service/env"
	"auth-service/service"
	"context"
	"github.com/danielgtaylor/huma/v2"
	"net/http"
)

func authController(api huma.API, options *env.Options) {
	huma.Register(api, huma.Operation{
		OperationID: "register",
		Method:      http.MethodPost,
		Path:        "/auth/register",
		Summary:     "Register a new user",
		Tags:        []string{"Auth"},
	}, func(ctx context.Context, req *dto.RegisterRequest) (*dto.RegisterResponse, error) {
		err := service.RegisterNewUser(req)
		if err != nil {
			return nil, err
		}
		return &dto.RegisterResponse{Status: http.StatusCreated}, nil
	})

	huma.Register(api, huma.Operation{
		OperationID: "login",
		Method:      http.MethodPost,
		Path:        "/auth/login",
		Summary:     "Login with username/email and password",
		Tags:        []string{"Auth"},
	}, func(ctx context.Context, req *dto.LoginRequest) (*dto.LoginResponse, error) {
		return service.LoginUser(req)
	})

	huma.Register(api, huma.Operation{
		OperationID: "refresh",
		Method:      http.MethodPost,
		Path:        "/auth/refresh",
		Summary:     "Refresh access token using a refresh token",
		Tags:        []string{"Auth"},
	}, func(ctx context.Context, req *dto.RefreshRequest) (*dto.LoginResponse, error) {
		return service.RefreshSession(req)
	})
}
