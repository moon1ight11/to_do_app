package handlers

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"log"
	"net/http"
	"regexp"
	"strings"
	"todoapp/internal/services"
	"todoapp/internal/storage/repos/users"
)

type UserHandler struct {
	userService *services.UserService
}

func NewUserHandler(userService *services.UserService) *UserHandler {
	return &UserHandler{userService: userService}
}

// добаление нового пользователя в БД
func (u *UserHandler) AddUser(c *gin.Context) {
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
		c.JSON((http.StatusInternalServerError), gin.H{"error": err.Error()})
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
		c.JSON((http.StatusInternalServerError), gin.H{"error": err.Error()})
		return
	}

	// если нет - отклоняем
	if emailExist {
		log.Println("Email already exist")
		c.JSON((http.StatusConflict), gin.H{"error": "Email already exist"})
		return
	}

	// добавляем пользователя в БД
	ID, err := u.userService.CreateUser(NewUser)
	if err != nil {
		log.Println(err)
		c.JSON((http.StatusInternalServerError), gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"user_id": ID})
}

// проверка существующего пользователя
func (u *UserHandler) CheckUser(c *gin.Context) {
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
		c.JSON((http.StatusBadRequest), gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"user": foundUser})
}

// обновление параметров пользователя
func (u *UserHandler) UpdateUser(c *gin.Context) {
	var UpdatedUser struct {
		User_id    uuid.UUID `json:"user_id"`
		User_name  *string   `json:"user_name"`
		User_pass  *string   `json:"user_pass"`
		User_email *string   `json:"user_email"`
	}

	// user_id получим из куков
	UpdatedUser.User_id, _ = uuid.Parse("8b1bbae9-6e4d-41dc-984f-3a4a6c0abb17")

	// получаем обновленного пользователя с фронта
	if err := c.ShouldBindJSON(&UpdatedUser); err != nil {
		log.Println("Error in ShouldBindJSON", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	// если обновляется почта
	if UpdatedUser.User_email != nil {
		// проверяем, похожа ли новая почта на почту
		pattern := `^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`
		matched, err := regexp.MatchString(pattern, *UpdatedUser.User_email)
		if err != nil {
			log.Println("Error in MatchString", err)
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// если нет - отклоняем
		if !matched {
			log.Println("New email not looks like email")
			c.JSON(http.StatusBadRequest, gin.H{"error": "New email not looks like email"})
			return
		}
	}

	// обновляем нужные поля
	err := u.userService.UpdateUser(UpdatedUser.User_name, UpdatedUser.User_pass, UpdatedUser.User_email, UpdatedUser.User_id)
	if err != nil {
		log.Println(err)
		c.JSON((http.StatusInternalServerError), gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"UpdatedUser": UpdatedUser})
}

// удаление пользователя
func (u *UserHandler) DeleteUser(c *gin.Context) {
	// получаем id пользователя из куков
	id, _ := uuid.Parse("387e408f-080f-4323-9d92-edaab6bf37a9")

	// проверяем существование пользователя
	_, err := u.userService.CheckUserByID(id)
	if err != nil {
		log.Println(err)
		c.JSON((http.StatusInternalServerError), gin.H{"error": err.Error()})
		return
	}

	// удаляем пользователя
	err = u.userService.DeleteUser(id)
	if err != nil {
		log.Println(err)
		c.JSON((http.StatusInternalServerError), gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "delete is complete"})
}
