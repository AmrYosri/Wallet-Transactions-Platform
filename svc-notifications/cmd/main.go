package main

import (
	"fmt"
	"net/http"
	"os"

	"svc-notifications/api/rest"
	"svc-notifications/external/mongodb"
	"svc-notifications/internal/notification"
	"svc-notifications/util/logger"

	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load("../.env")
	if err != nil {
		fmt.Println("Error loading .env file:", err)
	}

	port := os.Getenv("SERVER_PORT")
	fmt.Println("Server port is:", port)

	logger.Init()
	logger.Log.Info().Msg("Logger initialized")

	db, err := mongodb.Connect()
	if err != nil {
		logger.Log.Fatal().Err(err).Msg("Failed to connect to MongoDB")
	}
	logger.Log.Info().Msg("Connected to MongoDB")

	repo := notification.NewRepository(db)
	service := notification.NewService(repo)
	controller := rest.NewController(service)
	routes := rest.NewRouter(controller)

	logger.Log.Info().Msg("Starting server on port " + port)
	err = http.ListenAndServe(":"+port, routes)
	if err != nil {
		logger.Log.Fatal().Err(err).Msg("Server failed to start")
	}
}