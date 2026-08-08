package main

import (
	"fmt"
	"net/http"
	"os"

	"svc-user/api/rest"
	"svc-user/internal/user"
	"svc-user/util/logger"

	"svc-user/external/mongodb"
	"github.com/joho/godotenv"
)

func main() {


	err := godotenv.Load("../.env")
	if err != nil{
		fmt.Println("Error loading .env file :",err)
		
	}
	port := os.Getenv("SERVER_PORT")
	fmt.Println("Server port is :" ,port)

	logger.Init()
	logger.Log.Info().Msg("logger initialized")
	db,err := mongodb.Connect()

	if err != nil {
		logger.Log.Fatal().Err(err).Msg("failed to connect to MongoDB")

	}
	logger.Log.Info().Msg("connected to MongoDB")

	repo := user.NewRepository(db)
	service :=user.NewService(repo)
	controller :=rest.NewController(service)
	routes :=rest.NewRouter(controller)
	

	logger.Log.Info().Msg("Starting server on port " + port)


	err = http.ListenAndServe(":"+port,routes)
	if err != nil{
		logger.Log.Fatal().Err(err).Msg("server failed to start")
	}
} 