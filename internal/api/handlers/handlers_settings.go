package handlers

import (
	"log"
	"net/http"
	"time"
	"todoapp/internal/services"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type SettingsHandler struct {
	settingsService *services.SettingsService
}

func NewSettingsHandler(settingsService *services.SettingsService) *SettingsHandler {
	if settingsService == nil {
		panic("settings service cannot be nil")
	}
	return &SettingsHandler{settingsService: settingsService}
}

// получение настроек пользователя
func (s *SettingsHandler) GetSettings(c *gin.Context) {
	// получаем id
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		log.Println("Error in parse uuid")
		c.JSON((http.StatusBadRequest), gin.H{"error": "Invalid request body"})
		return
	}

	// находим настройки по id
	user_settings, err := s.settingsService.GetSettings(id)
	if err != nil {
		log.Println(err)
		c.JSON((http.StatusInternalServerError), gin.H{"error": err})
		return
	}

	c.JSON((http.StatusOK), gin.H{"settings": user_settings})
}

// изменение настроек пользователя
func (s *SettingsHandler) UpdateSettings(c *gin.Context) {
	var UpdatedSettings struct {
		User_id          uuid.UUID      `json:"id"`
		Default_duration *time.Duration `json:"default_duration"`
		Default_tz       *string        `json:"default_tz"`
	}

	// получаем настройки с фронта
	if err := c.ShouldBindJSON(&UpdatedSettings); err != nil {
		log.Println("Error in ShouldBindJSON")
		c.JSON((http.StatusBadRequest), gin.H{"error": "Invalid request body"})
		return
	}

	// изменяем настройки
	err := s.settingsService.UpdatedSettings(UpdatedSettings.User_id, UpdatedSettings.Default_duration, UpdatedSettings.Default_tz)
	if err != nil {
		log.Println(err)
		c.JSON((http.StatusInternalServerError), gin.H{"error": err})
		return
	}

	c.JSON((http.StatusOK), gin.H{"updated_settings": UpdatedSettings})
}
