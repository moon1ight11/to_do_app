package api

import (
	"todoapp/internal/api/routes"
	"github.com/gin-gonic/gin"
)

func Router () {
	userHandler := &routes.UserHandler{} 

	c := gin.Default()

	// регистрация
	c.POST("/v1/register", userHandler.AddUser)

	// авторизация
	c.GET("/v1/auth", userHandler.CheckUser)

	// обновление пользователя
	c.PATCH("v1/private/users", userHandler.UpdateUser)

	// удаление пользователя
	c.DELETE("v1/private/users", userHandler.DeleteUser)

	c.Run(":8080")
}