package service

import (
	"context"

	db "github.com/NatdanaiKhe/simplebank/db/sqlc"
	"github.com/NatdanaiKhe/simplebank/util"
)

type CreateUserParams struct {
	Username string
	Password string
	FullName string
	Email    string
}

type userService struct {
	store db.Store
}

type UserService interface {
	CreateUser(ctx context.Context, arg CreateUserParams) (db.User, error)
}

func NewUserService(store db.Store) UserService {
	return &userService{store: store}
}

func (s *userService) CreateUser(ctx context.Context, arg CreateUserParams) (db.User, error) {
	hashedPassword, err := util.HashPassword(arg.Password)
	if err != nil {
		return db.User{}, err
	}

	param := db.CreateUserParams{
		Username:       arg.Username,
		HashedPassword: hashedPassword,
		FullName:       arg.FullName,
		Email:          arg.Email,
	}
	user, err := s.store.CreateUser(ctx, param)
	return user, err
}
