package main

import (
	"fmt"
	"os"
	"strings"

	"go-echo-api/controllers"
	"go-echo-api/db"
	"go-echo-api/services"
	"go-echo-api/utils"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/labstack/gommon/log"
	"gopkg.in/go-playground/validator.v9"

	_ "go-echo-api/docs" // docs is generated by Swag CLI, you have to import it.

	"github.com/joho/godotenv"
	echoSwagger "github.com/swaggo/echo-swagger" // echo-swagger middleware
)

type Validator struct {
	validator *validator.Validate
}

func (v *Validator) Validate(i interface{}) error {
	return v.validator.Struct(i)
}

// @title Swagger Example API
// @version 1.0
// @description Conduit API
// @title Conduit API

// @host 127.0.0.1:3001
// @BasePath /api

// @schemes http https
// @produce	application/json
// @consumes application/json

// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	e := echo.New()
	e.Logger.SetLevel(log.DEBUG)
	e.Pre(middleware.RemoveTrailingSlash())
	e.Use(middleware.Logger())
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"},
		AllowHeaders: []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept, echo.HeaderAuthorization},
		AllowMethods: []string{echo.GET, echo.HEAD, echo.PUT, echo.PATCH, echo.POST, echo.DELETE},
	}))
	e.Use(middleware.GzipWithConfig(middleware.GzipConfig{
		Skipper: func(c echo.Context) bool {
			return strings.Contains(c.Request().URL.Path, "docs")
		},
	}))
	e.Validator = &Validator{
		validator: validator.New(),
	}

	e.GET("/docs/*", echoSwagger.WrapHandler)

	v1 := e.Group("/api")
	utils.GetClient()

	d := db.New()
	db.AutoMigrate(d)

	userService := services.NewUserService(d)
	articleService := services.NewArticleService(d)

	controllers.NewUserController(v1, userService)
	controllers.NewArticleController(v1, articleService, userService)
	fmt.Printf("Swagger UI is available on http://localhost:%s/docs/index.html", os.Getenv("HOST_PORT"))

	e.Logger.Fatal(e.Start(fmt.Sprintf("0.0.0.0:%s", os.Getenv("HOST_PORT"))))
}
