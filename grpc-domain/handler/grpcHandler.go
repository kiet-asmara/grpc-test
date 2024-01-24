package handler

import (
	"context"
	repository "ngc-grpc/grpc-domain/repo"
	"ngc-grpc/model"
)

type UserHandler struct {
	UserRepository repository.UserRepository
	model.UnimplementedUserServiceServer
}

func NewUserHandler(UserRepository repository.UserRepository) *UserHandler {
	return &UserHandler{
		UserRepository: UserRepository,
	}
}

func (h *UserHandler) AddUser(ctx context.Context, in *model.User) (*model.UserResponse, error) {
	var s = model.UserModel{
		ID:   in.Id,
		Name: in.Name,
	}

	_, err := h.UserRepository.Create(&s)
	if err != nil {
		return nil, err
	}

	var response = &model.UserResponse{
		Status: "Berhasil",
		Id:     in.Id,
		Name:   in.Name,
	}

	return response, nil
}

func (t *UserHandler) GetUserList(ctx context.Context, in *model.Empty) (*model.UserList, error) {
	users, err := t.UserRepository.ReadAll()
	if err != nil {
		return nil, err
	}

	var list []*model.User

	for _, v := range users {
		u := model.User{
			Id:   v.ID,
			Name: v.Name,
		}
		list = append(list, &u)
	}

	var response = &model.UserList{
		List: list,
	}

	return response, nil
}
