package models

import "github.com/google/uuid"

// модель для аутентификации
type UserAuth struct {
	Id    uuid.UUID `json:"user_id"`
	Name  string    `json:"user_name"`
	Email string    `json:"user_email" binding:"required,email"`
	Pass  string    `json:"user_pass"`
}

// модель для ответа
type UserRequest struct {
	Id    uuid.UUID `json:"user_id"`
	Name  string    `json:"user_name"`
	Email string    `json:"user_email" binding:"required,email"`
}

// модель для обновления
type UserUpdate struct {
	Id    uuid.UUID `json:"user_id"`
	Name  *string   `json:"user_name"`
	Email *string   `json:"user_email"`
	Pass  *string   `json:"user_pass"`
}
