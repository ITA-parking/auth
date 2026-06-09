package api

import (
	"auth-service/env"
	"github.com/danielgtaylor/huma/v2"
)

var controllers []Controller = []Controller{authController}

type Controller func(api huma.API, options *env.Options)

func RegisterAll(api huma.API, options *env.Options) {
	for _, c := range controllers {
		c(api, options)
	}
}
