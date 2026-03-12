package settingsservice

import (
	"todoapp/internal/storage/repos/settingsrepos"
)

type SettingsService struct {
	settingsRepo *settingsrepos.Repo
}

func NewSettingsService(settingsRepo *settingsrepos.Repo) *SettingsService {
	return &SettingsService{settingsRepo: settingsRepo}
}
