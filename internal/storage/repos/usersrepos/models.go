package usersrepos

import (
	"github.com/google/uuid"
)

type User struct {
	Id    uuid.UUID
	Name  string
	Email string
	Pass  string
}
