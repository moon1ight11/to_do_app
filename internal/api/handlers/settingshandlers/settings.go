package settingshandlers

import (
	"todoapp/internal/api/jwt"
	"todoapp/internal/services/settingsservice"
)

type SettingsHandler struct {
	settingsService *settingsservice.SettingsService
	jwtService      jwt.TokenService
}

func NewSettingsHandler(settingsService *settingsservice.SettingsService, jwtService jwt.TokenService) *SettingsHandler {
	return &SettingsHandler{
		settingsService: settingsService,
		jwtService:      jwtService}
}
