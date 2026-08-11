package settingsservice

import (
	"context"
	"fmt"
	"regexp"

	"github.com/google/uuid"
	"todoapp/internal/api/models"
)

func (s *SettingsService) GetSettings(ctx context.Context, userId uuid.UUID) (models.Setting, error) {
	ctx, span := s.tracer.Start(ctx, "service.GetSettings")
	defer span.End()

	settings, err := s.settingsRepo.SettingsById(ctx, userId)
	if err != nil {
		span.RecordError(err)
		return models.Setting{}, fmt.Errorf("settingsservice.GetSettings: %w", err)
	}

	return models.Setting{
		TimeDuration: settings.TimeDuration,
		UserTz:       settings.UserTz,
	}, nil
}

func (s *SettingsService) UpdateSettings(ctx context.Context, userId uuid.UUID, duration *float64, tz *string) error {
	ctx, span := s.tracer.Start(ctx, "service.UpdateSettings")
	defer span.End()

	if tz != nil {
		pattern := `^UTC([+-](?:1[0-4]|[0-9])(?::?[0-5][0-9])?)?$`
		matched, err := regexp.MatchString(pattern, *tz)
		if err != nil {
			span.RecordError(err)
			return fmt.Errorf("settingsservice.UpdateSettings: validate tz: %w", err)
		}
		if !matched {
			err := fmt.Errorf("timezone not valid: %s", *tz)
			span.RecordError(err)
			return fmt.Errorf("settingsservice.UpdateSettings: %w", err)
		}
	}

	transaction, err := s.settingsRepo.DB().BeginTx(ctx, nil)
	if err != nil {
		span.RecordError(err)
		return fmt.Errorf("settingsservice.UpdateSettings: begin tx: %w", err)
	}
	defer transaction.Rollback()

	if duration != nil {
		if err := s.settingsRepo.UpdateDuration(ctx, *duration, userId, transaction); err != nil {
			span.RecordError(err)
			return fmt.Errorf("settingsservice.UpdateSettings: update duration: %w", err)
		}
	}

	if tz != nil {
		if err := s.settingsRepo.UpdateTZ(ctx, *tz, userId, transaction); err != nil {
			span.RecordError(err)
			return fmt.Errorf("settingsservice.UpdateSettings: update tz: %w", err)
		}
	}

	if err := transaction.Commit(); err != nil {
		span.RecordError(err)
		return fmt.Errorf("settingsservice.UpdateSettings: commit: %w", err)
	}

	return nil
}
