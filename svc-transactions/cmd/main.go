package main

import (
	"fmt"
	"net/http"
	"os"

	"svc-transactions/api/rest"
	"svc-transactions/client/wallet"
	"svc-transactions/external/mongodb"
	"svc-transactions/internal/transactions"
	"svc-transactions/util/logger"

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

	walletServiceURL := os.Getenv("WALLET_SERVICE_URL")
	walletClient := wallet.NewClient(walletServiceURL)

	repo := transactions.NewRepository(db)
	service := transactions.NewService(repo, walletClient)
	controller := rest.NewController(service)
	routes := rest.NewRouter(controller)

	logger.Log.Info().Msg("Starting server on port " + port)
	err = http.ListenAndServe(":"+port, routes)
	if err != nil {
		logger.Log.Fatal().Err(err).Msg("Server failed to start")
	}
}
