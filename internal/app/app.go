package app

import (
	"github.com/labstack/echo/v4/middleware"
	"log"
	"strconv"
	"time"

	handler "github.com/MAVIKE/yad-backend/internal/delivery/http"
	"github.com/MAVIKE/yad-backend/internal/repository"
	"github.com/MAVIKE/yad-backend/internal/service"
	"github.com/MAVIKE/yad-backend/pkg/auth"
	"github.com/labstack/echo/v4"
	"github.com/spf13/viper"
)

func Run(configPath string) {
	if err := initConfig(configPath); err != nil {
		log.Fatalf("error initializing configs: %s", err.Error())
	}

	dbPrefix := viper.GetString("db.name") + "."

	db, err := repository.NewPostgresDB(repository.Config{
		Host:     viper.GetString(dbPrefix + "host"),
		Port:     viper.GetString(dbPrefix + "port"),
		Username: viper.GetString(dbPrefix + "username"),
		DBName:   viper.GetString(dbPrefix + "dbname"),
		SSLMode:  viper.GetString(dbPrefix + "sslmode"),
		Password: viper.GetString(dbPrefix + "password"),
	})

	if err != nil {
		log.Fatalf("failed to initialize db: %s", err.Error())
	}

	repos := repository.NewRepository(db)

	signingKey := viper.GetString("token.signing_key")
	if signingKey == "" {
		log.Fatalf("failed to get token signing key")
	}

	accessTokenTTL, err := strconv.Atoi(viper.GetString("token.access_token_ttl"))
	if err != nil || accessTokenTTL == 0 {
		log.Fatalf("failed to get access token TTL: %s", err.Error())
	}

	tokenManager, err := auth.NewManager(signingKey)
	if err != nil {
		log.Fatalf(err.Error())
	}

	deps := service.Deps{
		Repos:          repos,
		TokenManager:   tokenManager,
		AccessTokenTTL: time.Duration(accessTokenTTL) * time.Hour,
	}

	services := service.NewService(deps)
	handlers := handler.NewHandler(services, tokenManager)

	app := echo.New()
	app.Use(middleware.Logger())
	handlers.Init(app)

	if err := app.Start(viper.GetString("port")); err != nil {
		log.Fatalf("failed to listen: %s", err.Error())
	}
}

func initConfig(configPath string) error {
	viper.AddConfigPath(configPath)
	viper.SetConfigName("config")
	return viper.ReadInConfig()
}
