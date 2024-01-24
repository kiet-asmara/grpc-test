package main

import (
	"log"
	"ngc-grpc/model"
	"ngc-grpc/rest-domain/handler"

	"github.com/labstack/echo/v4"
	"google.golang.org/grpc"
)

func main() {
	connection, err := grpc.Dial(":50001", grpc.WithInsecure())
	if err != nil {
		log.Println(err)
	}

	serviceClient := model.NewUserServiceClient(connection)

	restHandler := handler.NewHandler(serviceClient)

	e := echo.New()
	e.POST("/users", restHandler.CreateUser)
	e.GET("/users", restHandler.GetUserList)

	err = e.Start(":8080")
	if err != nil {
		panic(err)
	}
}
