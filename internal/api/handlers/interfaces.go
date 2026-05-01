package handlers

import "github.com/gin-gonic/gin"

type AuthHandlerInterface interface {
	SignUp(c *gin.Context)
	SignIn(c *gin.Context)
	SignOut(c *gin.Context)
}

type UserHandlerInterface interface {
	GetUser(c *gin.Context)
	UpdateUser(*gin.Context)
	DeleteUser(*gin.Context)
}

type SettingsHandlerInterface interface {
	GetSettings(*gin.Context)
	UpdateSettings(*gin.Context)
}

type TaskHandlerInterface interface {
	CreateTask(*gin.Context)
	GetTasks(*gin.Context)
	GetTaskById(*gin.Context)
	UpdateTask(*gin.Context)
	DeleteTask(*gin.Context)
}
