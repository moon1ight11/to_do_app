package app

import (
    "todoapp/internal/storage"
)

type App struct {
	DBex *storage.DataBase
}