package models

// модель для обновления
type UpdatedSettings struct {
	TimeDuration *float64 `json:"time_duration"`
	UserTz       *string  `json:"user_tz" binding:"required"`
}

// модель для получения
type Setting struct {
	TimeDuration float64 `json:"time_duration"`
	UserTz       string  `json:"user_tz" binding:"required"`
}
