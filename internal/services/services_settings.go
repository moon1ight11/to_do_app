package services

import (
	"github.com/google/uuid"
	"todoapp/internal/storage/repos/settings"
)

type SettingsService struct {
	settingsRepo *settings.Base
}

func NewSettingsService(settingsRepo *settings.Base) *SettingsService {
	return &SettingsService{settingsRepo: settingsRepo}
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
func (s *SettingsService) UpdatedSettings(id uuid.UUID, new_duration *float64, new_location *string) error {
	// открываем транзакцию
	transaction, err := s.settingsRepo.DB.Begin()
	if err != nil {
		return err
	}

	// отложенно откатываем транзакцию
	defer transaction.Rollback()

	// если меняем продолжительность
	if new_duration != nil {
		err := s.settingsRepo.UpdateDuration(*new_duration, id, transaction)
		if err != nil {
			return err
		}
	}

	// если меняем таймзону
	if new_location != nil {
		err := s.settingsRepo.UpdateTZ(*new_location, id, transaction)
		if err != nil {
			return err
		}
	}

	// если все ок - подтверждаем транзакцию
	transaction.Commit()

	return nil
}
