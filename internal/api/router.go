package api

import (
	"todoapp/internal/api/routes"
	"github.com/gin-gonic/gin"
)

func Router () {
	userHandler := &routes.UserHandler{} 

	c := gin.Default()

	c.POST("/register", userHandler.AddUser)
	c.GET("/auth", userHandler.CheckUser)
	c.PATCH("/user_update", userHandler.UpdateUser)
	c.DELETE("/delete_user/:id", userHandler.DeleteUser)

	c.Run(":8080")
}