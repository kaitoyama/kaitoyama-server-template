package main

import (
	"database/sql"
	"log"
	"os"

	"github.com/getkin/kin-openapi/openapi3"
	"github.com/kaitoyama/kaitoyama-server-template/internal/infrastructure/db"
	handler "github.com/kaitoyama/kaitoyama-server-template/internal/interface"
	"github.com/kaitoyama/kaitoyama-server-template/openapi"
	"github.com/labstack/echo/v4"
)

func SetupRouter(database *sql.DB) *echo.Echo {
	// Initialize Echo
	e := echo.New()

	// Load OpenAPI spec from file
	swaggerBytes, err := os.ReadFile("openapi/swagger.yml")
	if err != nil {
		log.Fatalf("Failed to read swagger file: %v", err)
	}
	loader := openapi3.NewLoader()
	swagger, err := loader.LoadFromData(swaggerBytes)
	if err != nil {
		log.Fatalf("Failed to load swagger spec: %v", err)
	}
	if err = swagger.Validate(loader.Context); err != nil {
		log.Fatalf("Swagger spec validation error: %v", err)
	}

	// Apply middleware to validate incoming requests against the OpenAPI spec

	dbChecker := db.NewDBHealthChecker(database)
	healthHandler := handler.NewHandler(dbChecker)

	// oapi-codegenのRegisterHandlersを使用
	openapi.RegisterHandlers(e, healthHandler)

	return e
}
