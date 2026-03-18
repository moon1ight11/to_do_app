package api

import (
	"github.com/gin-gonic/gin"
	"todoapp/internal/api/handlers"
	"todoapp/internal/api/jwt"
	"todoapp/internal/api/middleware"
	"todoapp/pkg/logger"
)

type Router struct {
	userHandler     handlers.UserHandlerInterface
	settingsHandler handlers.SettingsHandlerInterface
	tasksHandler    handlers.TaskHandlerInterface
	authHandler     handlers.AuthHandlerInterface
	ginEngine       *gin.Engine
}

func NewRouter(
	userHandler handlers.UserHandlerInterface,
	settingsHandler handlers.SettingsHandlerInterface,
	tasksHandler handlers.TaskHandlerInterface,
	authHandler handlers.AuthHandlerInterface,
) *Router {
	return &Router{
		userHandler:     userHandler,
		settingsHandler: settingsHandler,
		tasksHandler:    tasksHandler,
		authHandler:     authHandler,
		ginEngine:       gin.Default(),
	}
}

func (r *Router) Init(jwtService jwt.TokenService, logger *logger.Logger) {
	r.ginEngine.Use(middleware.CORS())

	// группировка роутов
	privateGroup := r.ginEngine.Group("/v1/private")
	authGroup := r.ginEngine.Group("/v1/auth")

	// MIDDLEWARES //
	privateGroup.Use(middleware.Auth(jwtService, logger))

	// АУТЕНТИФИКАЦИЯ //
	// регистрация
	authGroup.POST("/sign-up", r.authHandler.SignUp)
	// авторизация
	authGroup.POST("/sign-in", r.authHandler.SignIn)
	// разлогин
	authGroup.POST("/sign-out", r.authHandler.SignOut)

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

func (r *Router) GetEngine() *gin.Engine {
	return r.ginEngine
}
