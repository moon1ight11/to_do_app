package api

import (
	"todoapp/internal/api/handlers"
	"github.com/gin-gonic/gin"
)

func UpRouter(userHandler *handlers.UserHandler) {

	defaultGroup := gin.Default()
	privateGroup := defaultGroup.Group("/v1/private")

	// регистрация
	defaultGroup.POST("/v1/register", userHandler.AddUser)

	// авторизация
	defaultGroup.GET("/v1/auth", userHandler.CheckUser)

	// обновление пользователя
	privateGroup.PATCH("/users", userHandler.UpdateUser)

	// удаление пользователя
	privateGroup.DELETE("/users", userHandler.DeleteUser)

	defaultGroup.Run(":8080")
}
