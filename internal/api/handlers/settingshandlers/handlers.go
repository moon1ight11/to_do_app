package settingshandlers

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"log"
	"net/http"
	"regexp"
	"todoapp/internal/api/models"
)

// получение настроек пользователя
func (s *SettingsHandler) GetSettings(c *gin.Context) {
	// получаем id из контекста
	userIDValue, exist := c.Get("UserId")
	if !exist {
		c.JSON(http.StatusForbidden, gin.H{"error": "User ID not found"})
		return
	}

	// приводим значение к uuid
	userId, ok := userIDValue.(uuid.UUID)
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID type"})
		return
	}

	// находим настройки по id
	settings, err := s.settingsService.GetSettings(userId)
	if err != nil {
		log.Println(err)
		c.JSON((http.StatusInternalServerError), gin.H{"error": err.Error()})
		return
	}

	c.JSON((http.StatusOK), gin.H{"settings": settings})
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
	userId, ok := userIDValue.(uuid.UUID)
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID type"})
		return
	}

	var updatedSettings models.UpdatedSettings

	// получаем настройки с фронта
	if err := c.ShouldBindJSON(&updatedSettings); err != nil {
		log.Println("Error in ShouldBindJSON")
		c.JSON((http.StatusBadRequest), gin.H{"error": "Invalid request body"})
		return
	}

	// если обновляется временная зона
	if updatedSettings.UserTz != nil {
		// проверяем, похожа ли входящая строа на временную зону
		pattern := `^UTC([+-](?:1[0-4]|[0-9])(?::?[0-5][0-9])?)?$`
		matched, err := regexp.MatchString(pattern, *updatedSettings.UserTz)
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
	err := s.settingsService.UpdateSettings(userId, updatedSettings.TimeDuration, updatedSettings.UserTz)
	if err != nil {
		log.Println(err)
		c.JSON((http.StatusInternalServerError), gin.H{"error": err.Error()})
		return
	}

	c.JSON((http.StatusOK), gin.H{"updated_settings": updatedSettings})
}
