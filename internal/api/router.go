package api

import (
	"net/http/pprof"

	"github.com/gin-gonic/gin"

	"todoapp/internal/api/handlers"
	"todoapp/internal/api/jwt"
	"todoapp/internal/api/middleware"
	"todoapp/internal/config"
	"todoapp/internal/metrics"
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

func (r *Router) Init(
	jwtService jwt.TokenService,
	l logger.LoggerInterface,
	cfg *config.Config,
	m *metrics.Metrics,
) {
	r.registerPprof()
	r.ginEngine.Use(middleware.Tracing(cfg.ServiceName))
	r.ginEngine.Use(m.GinMiddleware())
	r.ginEngine.Use(middleware.CORS())
	m.RegisterMetricsHandler(r.ginEngine, "/metrics")

	privateGroup := r.ginEngine.Group("/v1/private")
	authGroup := r.ginEngine.Group("/v1/auth")

	privateGroup.Use(middleware.Auth(jwtService, l))

	authGroup.POST("/sign-up", r.authHandler.SignUp)
	authGroup.POST("/sign-in", r.authHandler.SignIn)

	privateGroup.POST("/sign-out", r.authHandler.SignOut)

	privateGroup.GET("/users", r.userHandler.GetUser)
	privateGroup.PATCH("/users", r.userHandler.UpdateUser)
	privateGroup.DELETE("/users", r.userHandler.DeleteUser)

	privateGroup.GET("/settings", r.settingsHandler.GetSettings)
	privateGroup.PATCH("/settings", r.settingsHandler.UpdateSettings)

	privateGroup.POST("/tasks", r.tasksHandler.CreateTask)
	privateGroup.GET("/tasks", r.tasksHandler.GetTasks)
	privateGroup.GET("/tasks/:task_id", r.tasksHandler.GetTaskById)
	privateGroup.PATCH("/tasks/:task_id", r.tasksHandler.UpdateTask)
	privateGroup.DELETE("/tasks/:task_id", r.tasksHandler.DeleteTask)
}

func (r *Router) registerPprof() {
	r.ginEngine.GET("/debug/pprof/", gin.WrapF(pprof.Index))
	r.ginEngine.GET("/debug/pprof/cmdline", gin.WrapF(pprof.Cmdline))
	r.ginEngine.GET("/debug/pprof/profile", gin.WrapF(pprof.Profile))
	r.ginEngine.GET("/debug/pprof/symbol", gin.WrapF(pprof.Symbol))
	r.ginEngine.POST("/debug/pprof/symbol", gin.WrapF(pprof.Symbol))
	r.ginEngine.GET("/debug/pprof/trace", gin.WrapF(pprof.Trace))
	r.ginEngine.GET("/debug/pprof/heap", gin.WrapF(pprof.Index))
	r.ginEngine.GET("/debug/pprof/goroutine", gin.WrapF(pprof.Index))
	r.ginEngine.GET("/debug/pprof/threadcreate", gin.WrapF(pprof.Index))
	r.ginEngine.GET("/debug/pprof/block", gin.WrapF(pprof.Index))
	r.ginEngine.GET("/debug/pprof/mutex", gin.WrapF(pprof.Index))
	r.ginEngine.GET("/debug/pprof/allocs", gin.WrapF(pprof.Index))
}

func (r *Router) GetEngine() *gin.Engine {
	return r.ginEngine
}
