package settingsservice

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"regexp"
	"todoapp/internal/api/models"
)

// получение настроек
func (s *SettingsService) GetSettings(ctx context.Context, userId uuid.UUID) (models.Setting, error) {
	// получаем настройки из репозитория
	settings, err := s.settingsRepo.SettingsById(ctx, userId)
	if err != nil {
		return models.Setting{}, fmt.Errorf("error in GetSettings: %w", err)
	}

	// приводим тип
	var settingsApi models.Setting
	settingsApi.UserTz = settings.UserTz
	settingsApi.TimeDuration = settings.TimeDuration

	return settingsApi, nil
}

// изменение настроек
func (s *SettingsService) UpdateSettings(ctx context.Context, userId uuid.UUID, duration *float64, tz *string) error {
	// если обновляется временная зона
	if tz != nil {
		// проверяем, похожа ли входящая строа на временную зону
		pattern := `^UTC([+-](?:1[0-4]|[0-9])(?::?[0-5][0-9])?)?$`
		matched, err := regexp.MatchString(pattern, *tz)
		if err != nil {
			return fmt.Errorf("error in UpdateSettings: %w", err)
		}

		// если нет - прокидываем
		if !matched {
			return fmt.Errorf("error in UpdateSettings: timezone not looks like timezone")
		}
	}

	// открываем транзакцию
	transaction, err := s.settingsRepo.DB.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("error in UpdateSettings BeginTx: %w", err)
	}

	// отложенно откатываем транзакцию
	defer transaction.Rollback()

	// если меняем продолжительность
	if duration != nil {
		err := s.settingsRepo.UpdateDuration(ctx, *duration, userId, transaction)
		if err != nil {
			return fmt.Errorf("error in UpdateSettings: %w", err)
		}
	}

	// если меняем таймзону
	if tz != nil {
		err := s.settingsRepo.UpdateTZ(ctx, *tz, userId, transaction)
		if err != nil {
			return fmt.Errorf("error in UpdateSettings: %w", err)
		}
	}

	// если все ок - подтверждаем транзакцию
	err = transaction.Commit()
	if err != nil {
		return fmt.Errorf("error in UpdateSettings commit: %w", err)
	}

	return nil
}
