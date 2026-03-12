package settingsservice

import (
	"fmt"
	"log"
	"regexp"
	"todoapp/internal/api/models"

	"github.com/google/uuid"
)

// получение настроек
func (s *SettingsService) GetSettings(userId uuid.UUID) (models.Setting, error) {
	// получаем настройки из репозитория
	settings, err := s.settingsRepo.SettingsById(userId)
	if err != nil {
		return models.Setting{}, err
	}

	// приводим тип
	var settingsApi models.Setting
	settingsApi.UserTz = settings.UserTz
	settingsApi.TimeDuration = settings.TimeDuration

	return settingsApi, nil
}

// изменение настроек
func (s *SettingsService) UpdateSettings(userId uuid.UUID, duration *float64, tz *string) error {
	// если обновляется временная зона
	if tz != nil {
		// проверяем, похожа ли входящая строа на временную зону
		pattern := `^UTC([+-](?:1[0-4]|[0-9])(?::?[0-5][0-9])?)?$`
		matched, err := regexp.MatchString(pattern, *tz)
		if err != nil {
			log.Println("Error in MatchString", err)
			return fmt.Errorf("error in check email: %w", err)
		}

		// если нет - прокидываем
		if !matched {
			log.Println("Timezone not right")
			return fmt.Errorf("timezone not looks like timezone")
		}
	}

	// открываем транзакцию
	transaction, err := s.settingsRepo.DB.Begin()
	if err != nil {
		return err
	}

	// отложенно откатываем транзакцию
	defer transaction.Rollback()

	// если меняем продолжительность
	if duration != nil {
		err := s.settingsRepo.UpdateDuration(*duration, userId, transaction)
		if err != nil {
			return err
		}
	}

	// если меняем таймзону
	if tz != nil {
		err := s.settingsRepo.UpdateTZ(*tz, userId, transaction)
		if err != nil {
			return err
		}
	}

	// если все ок - подтверждаем транзакцию
	err = transaction.Commit()
	if err != nil {
		return fmt.Errorf("error in update settings commit: %w", err)
	}

	return nil
}
