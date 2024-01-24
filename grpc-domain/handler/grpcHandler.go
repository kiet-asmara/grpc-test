package handler

import (
	"context"
	"fmt"
	repository "ngc-grpc/grpc-domain/repo"
	"ngc-grpc/helpers"
	"ngc-grpc/model"

	"golang.org/x/crypto/bcrypt"
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

func (h *UserHandler) CreateUser(ctx context.Context, in *model.UserRegister) (*model.User, error) {
	// create hashed password
	hashedPass, err := bcrypt.GenerateFromPassword([]byte(in.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	var s = model.UserAll{
		ID:       in.Id,
		Name:     in.Name,
		Password: string(hashedPass),
	}

	_, err = h.UserRepository.Create(&s)
	if err != nil {
		return nil, err
	}

	var response = &model.User{
		Id:   in.Id,
		Name: in.Name,
	}

	return response, nil
}

func (h *UserHandler) VerifyUserCredentials(ctx context.Context, in *model.UserLogin) (*model.JWT, error) {
	var username = in.Name

	existingUser, err := h.UserRepository.ReadByName(username)
	if err != nil {
		return nil, err
	}

	// check password
	passCheck := bcrypt.CompareHashAndPassword([]byte(existingUser.Password), []byte(in.Password))
	if passCheck != nil {
		return nil, fmt.Errorf("invalid user or password")
	}

	// generate JWT
	token, err := helpers.GenerateJWT(existingUser.ID, existingUser.Name)
	if err != nil {
		return nil, err
	}

	var response = &model.JWT{
		Token: token,
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

func (t *UserHandler) GetUserByID(ctx context.Context, in *model.ID) (*model.User, error) {
	user, err := t.UserRepository.ReadID(in.Id)
	if err != nil {
		return nil, err
	}

	var response = &model.User{
		Id:   user.Id,
		Name: user.Name,
	}

	return response, nil
}

func (t *UserHandler) DeleteUser(ctx context.Context, in *model.ID) (*model.Empty, error) {
	res, err := t.UserRepository.Delete(in.Id)
	if err != nil {
		return nil, err
	}
	if res.DeletedCount == 0 {
		return nil, fmt.Errorf("user not found")
	}

	return &model.Empty{}, nil
}

func (t *UserHandler) UpdateUser(ctx context.Context, in *model.UserUpdate) (*model.User, error) {

	updateInput := model.UserModel{
		ID:   in.Newid,
		Name: in.Name,
	}

	res, err := t.UserRepository.Update(in.Id, &updateInput)
	if err != nil {
		return nil, err
	}
	if res.MatchedCount == 0 {
		return nil, fmt.Errorf("user not found")
	}

	var response = &model.User{
		Id:   updateInput.ID,
		Name: updateInput.Name,
	}

	return response, nil
}
