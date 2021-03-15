package app

import (
	handler "github.com/MAVIKE/yad-backend/internal/delivery/http"
	"github.com/MAVIKE/yad-backend/internal/repository"
	"github.com/MAVIKE/yad-backend/internal/service"
	"github.com/MAVIKE/yad-backend/pkg/auth"
	"github.com/labstack/echo/v4"
	"github.com/spf13/viper"
	"log"
)

const (
	DB_PREFIX = "local-db."
)

func Run(configPath string) {
	if err := initConfig(configPath); err != nil {
		log.Fatalf("error initializing configs: %s", err.Error())
	}

	db, err := repository.NewPostgresDB(repository.Config{
		Host:     viper.GetString(DB_PREFIX + "host"),
		Port:     viper.GetString(DB_PREFIX + "port"),
		Username: viper.GetString(DB_PREFIX + "username"),
		DBName:   viper.GetString(DB_PREFIX + "dbname"),
		SSLMode:  viper.GetString(DB_PREFIX + "sslmode"),
		Password: viper.GetString(DB_PREFIX + "password"),
	})

	if err != nil {
		log.Fatalf("failed to initialize db: %s", err.Error())
	}

	repos := repository.NewRepository(db)

	// TODO signing key from configs
	tokenManager, err := auth.NewManager(service.SIGNING_KEY)
	if err != nil {
		log.Fatalf(err.Error())
	}
	deps := service.Deps{
		Repos:          repos,
		TokenManager:   tokenManager,
		AccessTokenTTL: service.ACCESS_TOKEN_TTL,
	}

	services := service.NewService(deps)
	handlers := handler.NewHandler(services)

	app := echo.New()
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
