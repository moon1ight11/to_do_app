package usersservice

import "todoapp/internal/storage/repos/usersrepos"

type UserService struct {
	userRepo *usersrepos.Repo
}

func NewUserService(userRepo *usersrepos.Repo) *UserService {
	return &UserService{userRepo: userRepo}
}
