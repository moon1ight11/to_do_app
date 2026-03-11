package models

// обновление настроек
type UpdatedSettings struct {
	TimeDuration *float64 `json:"time_duration"`
	UserTz       *string  `json:"user_tz" binding:"required"`
}

// получение настроек
type Setting struct {
	TimeDuration float64 `json:"time_duration"`
	UserTz       string  `json:"user_tz" binding:"required"`
}
