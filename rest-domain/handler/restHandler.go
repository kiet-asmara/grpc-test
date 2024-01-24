package handler

import (
	"context"
	"net/http"
	"ngc-grpc/model"

	"github.com/labstack/echo/v4"
)

type Handler struct {
	Grpc model.UserServiceClient
}

func NewHandler(grpc model.UserServiceClient) *Handler {
	return &Handler{
		Grpc: grpc,
	}
}

func (h *Handler) CreateUser(c echo.Context) error {
	// bind json input
	var input model.UserModel
	if err := c.Bind(&input); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{
			"message": "invalid request" + err.Error(),
		})
	}

	in := model.User{
		Id:   input.ID,
		Name: input.Name,
	}

	response, err := h.Grpc.AddUser(context.Background(), &in)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{
			"message": err,
		})
	}

	return c.JSON(http.StatusCreated, response)
}

func (h *Handler) GetUserList(c echo.Context) error {
	list, err := h.Grpc.GetUserList(context.Background(), &model.Empty{})
	if err != nil {
		return c.JSON(500, echo.Map{
			"message": err,
		})
	}
	var result []*model.UserModel
	for _, v := range list.List {
		u := model.UserModel{
			ID:   v.Id,
			Name: v.Name,
		}
		result = append(result, &u)
	}

	return c.JSON(200, result)
}
