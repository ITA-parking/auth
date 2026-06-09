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
	//basePath := "/auth"

	//TODO: add session refresh endpoint!!!!
}
