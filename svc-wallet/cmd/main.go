package main

import (
	"fmt"
	"net/http"
	"os"
	"svc-wallet/util/logger"

	"svc-wallet/external/mongodb"

	"svc-wallet/internal/wallet"

	"svc-wallet/api/rest"

	"svc-wallet/client/user"

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
	userClient, err := user.NewClient("localhost:9002")
	if err != nil {
		logger.Log.Fatal().Err(err).Msg("Failed to connect to svc-user gRPC")
	}

	repo := wallet.NewRepository(db)
	service := wallet.NewService(repo, userClient)
	controller := rest.NewController(service)
	routes := rest.NewRouter(controller)

	logger.Log.Info().Msg("Starting server on port " + port)
	err = http.ListenAndServe(":"+port, routes)
	if err != nil {
		logger.Log.Fatal().Err(err).Msg("Server failed to start")
	}
}
