package handlers

import (
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"strings"
	"todoapp/internal/api/jwt"
	"todoapp/internal/services"
	"todoapp/internal/storage/repos/users"
)

type AuthHandler struct {
	userService *services.UserService
	jwtService  jwt.TokenService
}

func NewAuthHandler(userService *services.UserService, jwtService jwt.TokenService) *AuthHandler {
	return &AuthHandler{
		userService: userService,
		jwtService:  jwtService,
	}
}

// регистрация
func (u *AuthHandler) SignUp(c *gin.Context) {
	// получаем пользователя с фронта
	var NewUser users.User
	if err := c.ShouldBindJSON(&NewUser); err != nil {
		log.Println("Error in ShouldBindJSON", err)
		c.JSON((http.StatusBadRequest), gin.H{"error": err.Error()})
		return
	}

	// проверка валидности имени пользователя
	if strings.TrimSpace(NewUser.Name) == "" {
		log.Println("Name is empty")
		c.JSON((http.StatusBadRequest), gin.H{"error": "Name is empty"})
		return
	}

	// проверка валидности пароля
	if strings.TrimSpace(NewUser.Pass) == "" {
		log.Println("Pass is empty")
		c.JSON((http.StatusBadRequest), gin.H{"error": "Pass is empty"})
		return
	}

	// проверяем, свободно ли имя пользователя
	nameExist, err := u.userService.CheckName(NewUser.Name)
	if err != nil {
		log.Println(err)
		c.JSON((http.StatusConflict), gin.H{"error": err.Error()})
		return
	}

	// если нет - отклоняем
	if nameExist {
		log.Println("Name already exist")
		c.JSON((http.StatusConflict), gin.H{"error": "Name already exist"})
		return
	}

	// проверяем, свободна ли указанная почта
	emailExist, err := u.userService.CheckEmail(NewUser.Email)
	if err != nil {
		log.Println(err)
		c.JSON((http.StatusConflict), gin.H{"error": err.Error()})
		return
	}

	// если нет - отклоняем
	if emailExist {
		log.Println("Email already exist")
		c.JSON((http.StatusConflict), gin.H{"error": "Email already exist"})
		return
	}

	// добавляем пользователя в БД
	user_id, err := u.userService.AddUser(NewUser)
	if err != nil {
		log.Println(err)
		c.JSON((http.StatusInternalServerError), gin.H{"error": err.Error()})
		return
	}

	// генерируем токен для нового пользователя
	token, err := u.jwtService.GenerateToken(user_id, NewUser.Name, NewUser.Email)
	if err != nil {
		log.Println(err)
		c.JSON((http.StatusInternalServerError), gin.H{"error": err.Error()})
		return
	}

	// устанавливаем куки
	c.SetCookie("cookie", token, 3600, "/", "", false, true)

	c.JSON(http.StatusCreated, gin.H{"user_id": user_id})
}

// авторизация
func (u *AuthHandler) SignIn(c *gin.Context) {
	// получаем пользователя с фронта
	var user users.User
	if err := c.ShouldBindJSON(&user); err != nil {
		log.Println("Error in ShouldBindJSON", err)
		c.JSON((http.StatusBadRequest), gin.H{"error": err.Error()})
		return
	}

	// проверяем пароль на валидность
	if strings.TrimSpace(user.Pass) == "" {
		log.Println("Pass is empty")
		c.JSON((http.StatusBadRequest), gin.H{"error": "Pass is empty"})
		return
	}

	// проверка пользователя
	_, foundUser, err := u.userService.CheckUser(user)
	if err != nil {
		log.Println(err)
		c.JSON((http.StatusUnauthorized), gin.H{"error": err.Error()})
		return
	}

	// генерируем токен для найденного пользователя
	token, err := u.jwtService.GenerateToken(foundUser.Id, foundUser.Name, foundUser.Email)
	if err != nil {
		log.Println(err)
		c.JSON((http.StatusInternalServerError), gin.H{"error": err.Error()})
		return
	}

	// устанавливаем куки
	c.SetCookie("cookie", token, 3600, "/", "", false, true)

	c.JSON(http.StatusOK, gin.H{"user": foundUser})
}

// разлогин
func (u *AuthHandler) SignOut(c *gin.Context) {
	c.SetCookie("cookie", "1", -1, "/", "", false, false)

	c.JSON(http.StatusOK, gin.H{"message": "success"})
}
