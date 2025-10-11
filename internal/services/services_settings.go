package services

import (
	"time"
	"todoapp/internal/storage/repos/settings"
	"github.com/google/uuid"
)

type SettingsService struct {
	settingsRepo *settings.Base
}



// получение настроек
func (s *SettingsService) GetSettings(id uuid.UUID) (settings.Setting, error) {
	UserSettings, err := s.settingsRepo.SettingsById(id)
	if err != nil {
		return settings.Setting{}, err
	}

	return UserSettings, nil
}

// изменение настроек
func (s *SettingsService) UpdatedSettings(id uuid.UUID, new_duration *time.Duration, new_location *string) error {
	// если меняем продолжительность
	if new_duration != nil {
		err := s.settingsRepo.UpdateDuration(*new_duration, id)
		if err != nil {
			return err
		}
	}

	// если меняем таймзону
	if new_location != nil {
		err := s.settingsRepo.UpdateTZ(*new_location, id)
		if err != nil {
			return err
		}
	}

	return nil
}