package main

import (
	"fmt"
	"log"
	"net"
	"ngc-grpc/config"
	"ngc-grpc/grpc-domain/handler"
	repository "ngc-grpc/grpc-domain/repo"
	"ngc-grpc/model"

	"google.golang.org/grpc"
)

func main() {
	db, err := config.InitDB()
	if err != nil {
		log.Println(err)
	}

	mongoRepo := repository.NewMongoRepository(db)
	userHandler := handler.NewUserHandler(mongoRepo)

	grpcServer := grpc.NewServer()
	model.RegisterUserServiceServer(grpcServer, userHandler)

	// start gRPC server
	listen, err := net.Listen("tcp", ":50001")
	if err != nil {
		log.Println(err)
	}

	err = grpcServer.Serve(listen)
	if err != nil {
		log.Println(err)
	}

	fmt.Println("INI GRPC")
}
