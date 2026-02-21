package models

import "github.com/google/uuid"

type User struct {
	Id    uuid.UUID `json:"user_id"`
	Name  string    `json:"user_name"`
	Email string    `json:"user_email" binding:"required,email"`
	Pass  string    `json:"user_pass"`
}