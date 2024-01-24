package main

import (
	"log"
	"ngc-grpc/model"
	"ngc-grpc/rest-domain/handler"
	"ngc-grpc/rest-domain/middleware"

	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	"google.golang.org/grpc"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal(err)
	}

	connection, err := grpc.Dial(":50001", grpc.WithInsecure())
	if err != nil {
		log.Fatal(err)
	}

	serviceClient := model.NewUserServiceClient(connection)

	restHandler := handler.NewHandler(serviceClient)

	e := echo.New()

	e.POST("/users", restHandler.CreateUser)
	e.POST("/login", restHandler.LoginUser)
	g := e.Group("")
	g.Use(middleware.AuthMiddleware)
	g.GET("/users", restHandler.GetUserList)
	g.GET("/users/:id", restHandler.GetUserByID)
	g.DELETE("/users/:id", restHandler.DeleteUser)
	g.PUT("/users/:id", restHandler.UpdateUser)

	err = e.Start(":8080")
	if err != nil {
		panic(err)
	}
}
