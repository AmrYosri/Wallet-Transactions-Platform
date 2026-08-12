package main

import (
	"fmt"
	"net"
	"net/http"
	"os"

	grpcapi "svc-wallet/api/grpc"
	"svc-wallet/api/rest"
	"svc-wallet/client/user"
	"svc-wallet/external/mongodb"
	"svc-wallet/internal/wallet"
	"svc-wallet/proto"
	"svc-wallet/util/logger"

	"github.com/joho/godotenv"
	"google.golang.org/grpc"
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

	userClient, err := user.NewClient(os.Getenv("USER_SERVICE_GRPC_ADDR"))
	if err != nil {
		logger.Log.Fatal().Err(err).Msg("Failed to connect to svc-user gRPC")
	}

	repo := wallet.NewRepository(db)
	service := wallet.NewService(repo, userClient)

	grpcPort := os.Getenv("GRPC_PORT")
	grpcServer := grpc.NewServer()
	proto.RegisterWalletServiceServer(grpcServer, grpcapi.NewGRPCServer(service))

	listener, err := net.Listen("tcp", ":"+grpcPort)
	if err != nil {
		logger.Log.Fatal().Err(err).Msg("Failed to listen for gRPC")
	}

	go func() {
		logger.Log.Info().Msg("Starting gRPC server on port " + grpcPort)
		if err := grpcServer.Serve(listener); err != nil {
			logger.Log.Fatal().Err(err).Msg("gRPC server failed")
		}
	}()

	controller := rest.NewController(service)
	routes := rest.NewRouter(controller)

	logger.Log.Info().Msg("Starting server on port " + port)
	err = http.ListenAndServe(":"+port, routes)
	if err != nil {
		logger.Log.Fatal().Err(err).Msg("Server failed to start")
	}
}