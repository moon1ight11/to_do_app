package handlers

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"log"
	"net/http"
	"regexp"
	"todoapp/internal/services"
)

type SettingsHandler struct {
	settingsService *services.SettingsService
}

func NewSettingsHandler(settingsService *services.SettingsService) *SettingsHandler {
	return &SettingsHandler{settingsService: settingsService}
}

// получение настроек пользователя
func (s *SettingsHandler) GetSettings(c *gin.Context) {
	// получаем id из куков
	id, _ := uuid.Parse("8b1bbae9-6e4d-41dc-984f-3a4a6c0abb17")

	// находим настройки по id
	user_settings, err := s.settingsService.GetSettings(id)
	if err != nil {
		log.Println(err)
		c.JSON((http.StatusInternalServerError), gin.H{"error": err.Error()})
		return
	}

	c.JSON((http.StatusOK), gin.H{"settings": user_settings})
}

// изменение настроек пользователя
func (s *SettingsHandler) UpdateSettings(c *gin.Context) {
	// получаем id из куков
	id, _ := uuid.Parse("8b1bbae9-6e4d-41dc-984f-3a4a6c0abb17")

	var UpdatedSettings struct {
		Default_duration *float64 `json:"default_duration"`
		Default_tz       *string  `json:"default_tz"`
	}

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

		// если нет - опрокидываем
		if !matched {
			log.Println("Timezone not right")
			c.JSON(http.StatusBadRequest, gin.H{"error": "Timezone not right"})
			return
		}
	}

	// изменяем настройки
	err := s.settingsService.UpdatedSettings(id, UpdatedSettings.Default_duration, UpdatedSettings.Default_tz)
	if err != nil {
		log.Println(err)
		c.JSON((http.StatusInternalServerError), gin.H{"error": err.Error()})
		return
	}

	c.JSON((http.StatusOK), gin.H{"message": "settings updated succesfull"})
}
