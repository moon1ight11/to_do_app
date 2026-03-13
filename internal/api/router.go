package api

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"todoapp/internal/api/handlers/authhandlers"
	"todoapp/internal/api/handlers/settingshandlers"
	"todoapp/internal/api/handlers/taskshandlers"
	"todoapp/internal/api/handlers/usershandlers"
	"todoapp/internal/api/jwt"
	"todoapp/internal/api/middleware"
)

type Router struct {
	userHandler     *usershandlers.UserHandler
	settingsHandler *settingshandlers.SettingsHandler
	tasksHandler    *taskshandlers.TasksHandler
	authHandler     *authhandlers.AuthHandler
	ginEngine       *gin.Engine
}

func NewRouter(
	userHandler *usershandlers.UserHandler,
	settingsHandler *settingshandlers.SettingsHandler,
	tasksHandler *taskshandlers.TasksHandler,
	authHandler *authhandlers.AuthHandler,
) *Router {
	return &Router{
		userHandler:     userHandler,
		settingsHandler: settingsHandler,
		tasksHandler:    tasksHandler,
		authHandler:     authHandler,
		ginEngine:       gin.Default(),
	}
}

func (r *Router) Run() error {
	err := r.ginEngine.Run(":8080")
	if err != nil {
		return fmt.Errorf("error in run router: %w", err)
	}

	return nil
}

func (r *Router) Init(jwtService jwt.TokenService) {
	r.ginEngine.Use(middleware.CORS())

	// группировка роутов
	privateGroup := r.ginEngine.Group("/v1/private")
	authGroup := r.ginEngine.Group("/v1/auth")

	// MIDDLEWARES //
	privateGroup.Use(middleware.Auth(jwtService))

	// АУТЕНТИФИКАЦИЯ //
	// регистрация
	authGroup.POST("/sign-up", r.authHandler.SignUp)
	// авторизация
	authGroup.POST("/sign-in", r.authHandler.SignIn)
	// разлогин
	authGroup.GET("/sign-out", r.authHandler.SignOut)

	// ЮЗЕРЫ //
	// получение данных пользователя
	privateGroup.GET("/users", r.userHandler.GetUser)
	// обновление пользователя
	privateGroup.PATCH("/users", r.userHandler.UpdateUser)
	// удаление пользователя
	privateGroup.DELETE("/users", r.userHandler.DeleteUser)

	// НАСТРОЙКИ //
	// получение настроек пользователя
	privateGroup.GET("/settings", r.settingsHandler.GetSettings)
	// изменение настроек
	privateGroup.PATCH("/settings", r.settingsHandler.UpdateSettings)

	// ЗАДАЧИ //
	// создание задачи
	privateGroup.POST("/tasks", r.tasksHandler.CreateTask)
	// получение всех задач пользователя
	privateGroup.GET("/tasks", r.tasksHandler.GetTasks)
	// получение одной задачи по id
	privateGroup.GET("/tasks/:task_id", r.tasksHandler.GetTaskById)
	// изменение задачи
	privateGroup.PATCH("/tasks/:task_id", r.tasksHandler.UpdateTask)
	// удаление задачи
	privateGroup.DELETE("/tasks/:task_id", r.tasksHandler.DeleteTask)
}
