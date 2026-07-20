package main

import (
	"fmt"
	"os"
	"svc-wallet/util/logger"

	"svc-wallet/external/mongodb"

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

	_ , err = mongodb.Connect()
	if err != nil {
		logger.Log.Fatal().Err(err).Msg("Failed to connect to MongoDB")
	}
	logger.Log.Info().Msg("Connected to MongoDB")
}
