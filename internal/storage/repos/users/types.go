package users

import "todoapp/internal/storage"

type Base struct {
	storage.DataBase
}

type User struct {
	Id    int
	Name  string
	Email string
	Pass  string
}
