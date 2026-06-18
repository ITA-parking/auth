package main

import (
	"auth-service/api"
	"auth-service/env"
	"auth-service/logs"
	"auth-service/repo"
	"auth-service/service"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/danielgtaylor/huma/v2"
	"github.com/danielgtaylor/huma/v2/adapters/humachi"
	"github.com/danielgtaylor/huma/v2/humacli"
	"github.com/go-chi/chi/v5"

	_ "github.com/danielgtaylor/huma/v2/formats/cbor"
)

func main() {
	logs.InitLogger()

	cli := humacli.New(func(hooks humacli.Hooks, options *env.Options) {
		dbErr := repo.InitDB(options.DbConnectionString)
		if dbErr != nil {
			slog.Error("Failed to init DB", "error", dbErr)
			panic(dbErr)
		}
		// Create a new router & API
		router := chi.NewMux()
		router.Use(service.AuthMiddleware)
		router.Get("/health", func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			w.Write([]byte(`{"status":"ok","service":"auth"}`))
		})
		addDocs(router)
		apiHuma := humachi.New(router, huma.DefaultConfig("Auth API", "1.0.0"))
		api.RegisterAll(apiHuma, options)

		// Tell the CLI how to start your router.
		hooks.OnStart(func() {
			runningOn := fmt.Sprintf("0.0.0.0:%d", options.HttpPort)
			slog.Info("Server starting on http://" + runningOn)
			http.ListenAndServe(runningOn, router)
		})
	})

	cli.Run()
}

func addDocs(router *chi.Mux) {
	router.Get("/docs", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")
		w.Write([]byte(`<!doctype html>
						<html lang="en">
						  <head>
							<meta charset="utf-8" />
							<meta name="referrer" content="same-origin" />
							<meta name="viewport" content="width=device-width, initial-scale=1, shrink-to-fit=no" />
							<title>Docs Example reference</title>
							<!-- Embed elements Elements via Web Component -->
							<link href="https://unpkg.com/@stoplight/elements@8.0.0/styles.min.css" rel="stylesheet" />
							<script src="https://unpkg.com/@stoplight/elements@8.0.0/web-components.min.js"
									integrity="sha256-yIhuSFMJJ6mp2XTUAb4SiSYneP3Qav8Uu+7NBhGJW5A="
									crossorigin="anonymous"></script>
						  </head>
						  <body style="height: 100vh;">
							<elements-api
							  apiDescriptionUrl="/openapi.yaml"
							  router="hash"
							  layout="stacked"
							  tryItCredentialsPolicy="same-origin"
							/>
						  </body>
						</html>`))
	})
}
