package main

import (
	"database/sql"
	"embed"
	"io/fs"
	"net/http"
	"os"
	"path/filepath"

	"github.com/getkin/kin-openapi/openapi3"
	"github.com/kaitoyama/kaitoyama-server-template/internal/infrastructure/db"
	"github.com/kaitoyama/kaitoyama-server-template/internal/interface/handler"
	"github.com/kaitoyama/kaitoyama-server-template/openapi"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/rs/zerolog/log"
)

//go:embed static/*
var staticFiles embed.FS

//go:embed openapi/swagger.yml
var swaggerFile []byte

func SetupRouter(database *sql.DB) *echo.Echo {
	// Initialize Echo
	e := echo.New()

	// Add middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.CORS())

	// Setup static file serving
	setupStaticFileServing(e)

	// Load OpenAPI spec (use embedded file for buildpack compatibility)
	var swaggerBytes []byte
	var err error

	// Try to read from embedded file first (for buildpack)
	if len(swaggerFile) > 0 {
		swaggerBytes = swaggerFile
	} else {
		// Fallback to file system (for development)
		swaggerBytes, err = os.ReadFile("openapi/swagger.yml")
		if err != nil {
			log.Fatal().Err(err).Msg("Failed to read swagger file")
		}
	}

	loader := openapi3.NewLoader()
	swagger, err := loader.LoadFromData(swaggerBytes)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to load swagger spec")
	}
	if err = swagger.Validate(loader.Context); err != nil {
		log.Fatal().Err(err).Msg("Swagger spec validation error")
	}

	// Setup API routes
	api := e.Group("/api")
	
	dbChecker := db.NewDBHealthChecker(database)
	healthHandler := handler.NewHandler(dbChecker)

	// Register API handlers with /api prefix
	openapi.RegisterHandlersWithBaseURL(api, healthHandler, "")

	return e
}

// setupStaticFileServing configures static file serving for both development and production
func setupStaticFileServing(e *echo.Echo) {
	// Check if we're running with embedded files (production/buildpack)
	if staticFS, err := fs.Sub(staticFiles, "static"); err == nil {
		// Use embedded files
		e.GET("/static/*", echo.WrapHandler(http.FileServer(http.FS(staticFS))))
		log.Info().Msg("Using embedded static files")
	} else {
		// Development mode - serve from local filesystem
		staticPath := "static"
		if _, err := os.Stat(staticPath); err == nil {
			e.Static("/static", staticPath)
			log.Info().Msg("Using local static files from ./static")
		} else {
			log.Warn().Msg("Static directory not found, static file serving disabled")
		}
	}
}
