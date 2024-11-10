package main

import (
	"database/sql"
	"gateway-service/config"
	userHandler "gateway-service/handlers/users"
	"gateway-service/proto/cart"
	"gateway-service/repository/users"
	"gateway-service/routes"
	userSvc "gateway-service/usecases/users"

	cartHandler "gateway-service/handlers/cart"
	cartSvc "gateway-service/usecases/cart"

	"github.com/go-playground/validator"
	"github.com/go-redis/redis/v8"
	"google.golang.org/grpc"
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

	rpcConn, err := config.RpcDial(cfg.CartServicePort)
	if err != nil {
		return
	}

	validator := validator.New()

	routes := setupRoutes(dbConn, validator, redisConn, rpcConn)
	routes.Run(cfg.AppPort)
}

func setupRoutes(db *sql.DB, validator *validator.Validate, redis *redis.Client, grpc *grpc.ClientConn) *routes.Routes {
	userStore := users.NewStore(db)
	userSvc := userSvc.NewUserSvc(userStore, redis)
	userHandler := userHandler.NewHandler(userSvc, validator)

	cartInt := cart.NewCartServiceClient(grpc)
	cartSvc := cartSvc.NewSvc(cartInt)
	cartHandler := cartHandler.NewHandler(cartSvc)

	return &routes.Routes{
		User: userHandler,
		Cart: cartHandler,
	}
}
