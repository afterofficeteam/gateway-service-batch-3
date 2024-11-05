package main

import (
	"database/sql"
	"gateway-service/config"
	userHandler "gateway-service/handlers/users"
	"gateway-service/repository/users"
	"gateway-service/routes"
	userSvc "gateway-service/usecases/users"

	"github.com/go-playground/validator"
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

	validator := validator.New()

	routes := setupRoutes(dbConn, validator)
	routes.Run(cfg.AppPort)
}

func setupRoutes(db *sql.DB, validator *validator.Validate) *routes.Routes {
	userStore := users.NewStore(db)
	userSvc := userSvc.NewUserSvc(userStore)
	userHandler := userHandler.NewHandler(userSvc, validator)

	return &routes.Routes{
		User: userHandler,
	}
}
