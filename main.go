package main

import (
	"database/sql"
	"gateway-service/config"
	userHandler "gateway-service/handlers/users"
	"gateway-service/repository/users"
	"gateway-service/routes"
	userSvc "gateway-service/usecases/users"

	"github.com/go-playground/validator"
	"github.com/go-redis/redis/v8"
)

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		return
	}

	dbConn, err := config.ConnectToDatabase(config.Connection{
		Host:     cfg.DBHost,
		Port:     cfg.DBPort,
		User:     cfg.DBUser,
		Password: cfg.DBPassword,
		DBName:   cfg.DBName,
	})
	if err != nil {
		return
	}
	defer dbConn.Close()

	redisConn, err := config.ConnectToRedis(config.RedisConnection{
		Host:     cfg.RedisHost,
		Port:     cfg.RedisPort,
		Password: cfg.RedisPassword,
		DB:       cfg.RedisDB,
	})
	if err != nil {
		return
	}

	validator := validator.New()

	routes := setupRoutes(dbConn, validator, redisConn)
	routes.Run(cfg.AppPort)
}

func setupRoutes(db *sql.DB, validator *validator.Validate, redis *redis.Client) *routes.Routes {
	userStore := users.NewStore(db)
	userSvc := userSvc.NewUserSvc(userStore, redis)
	userHandler := userHandler.NewHandler(userSvc, validator)

	return &routes.Routes{
		User: userHandler,
	}
}
