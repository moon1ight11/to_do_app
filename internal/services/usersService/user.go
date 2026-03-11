package usersservice

import (
	"todoapp/internal/storage/repos/users"
)

type UserService struct {
	userRepo *users.Repo
}

func NewUserService(userRepo *users.Repo) *UserService {
	return &UserService{userRepo: userRepo}
}