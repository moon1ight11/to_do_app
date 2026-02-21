package handlers

import (
	"log"
	"net/http"
	"regexp"
	"todoapp/internal/api/jwt"
	"todoapp/internal/api/models"
	"todoapp/internal/services"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type SettingsHandler struct {
	settingsService *services.SettingsService
	jwtService      jwt.TokenService
}

func NewSettingsHandler(settingsService *services.SettingsService, jwtService jwt.TokenService) *SettingsHandler {
	return &SettingsHandler{
		settingsService: settingsService,
		jwtService:      jwtService}
}

// получение настроек пользователя
func (s *SettingsHandler) GetSettings(c *gin.Context) {
	// получаем id из контекста
	userIDValue, exist := c.Get("UserId")
	if !exist {
		c.JSON(http.StatusForbidden, gin.H{"error": "User ID not found"})
		return
	}

	// приводим значение к uuid
	UserId, ok := userIDValue.(uuid.UUID)
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID type"})
		return
	}

	// находим настройки по id
	user_settings, err := s.settingsService.GetSettings(UserId)
	if err != nil {
		log.Println(err)
		c.JSON((http.StatusInternalServerError), gin.H{"error": err.Error()})
		return
	}

	c.JSON((http.StatusOK), gin.H{"settings": user_settings})
}

// изменение настроек пользователя
func (s *SettingsHandler) UpdateSettings(c *gin.Context) {
	// получаем id из контекста
	userIDValue, exist := c.Get("UserId")
	if !exist {
		c.JSON(http.StatusForbidden, gin.H{"error": "User ID not found"})
		return
	}

	// приводим значение к uuid
	UserId, ok := userIDValue.(uuid.UUID)
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID type"})
		return
	}
	
	var UpdatedSettings models.UpdatedSettings

	// получаем настройки с фронта
	if err := c.ShouldBindJSON(&UpdatedSettings); err != nil {
		log.Println("Error in ShouldBindJSON")
		c.JSON((http.StatusBadRequest), gin.H{"error": "Invalid request body"})
		return
	}

	// если обновляется временная зона
	if UpdatedSettings.Default_tz != nil {
		// проверяем, похожа ли входящая строа на временную зону
		pattern := `^UTC([+-](?:1[0-4]|[0-9])(?::?[0-5][0-9])?)?$`
		matched, err := regexp.MatchString(pattern, *UpdatedSettings.Default_tz)
		if err != nil {
			log.Println("Error in MatchString", err)
			c.JSON(http.StatusBadRequest, gin.H{"error": err})
			return
		}

		// если нет - прокидываем
		if !matched {
			log.Println("Timezone not right")
			c.JSON(http.StatusBadRequest, gin.H{"error": "Timezone not right"})
			return
		}
	}

	// изменяем настройки
	err := s.settingsService.UpdateSettings(UserId, UpdatedSettings.Default_duration, UpdatedSettings.Default_tz)
	if err != nil {
		log.Println(err)
		c.JSON((http.StatusInternalServerError), gin.H{"error": err.Error()})
		return
	}

	c.JSON((http.StatusOK), gin.H{"updated_settings": UpdatedSettings})
}
