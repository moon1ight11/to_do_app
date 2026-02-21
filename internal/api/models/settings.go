package models

// обновление настроек
type UpdatedSettings struct {
	Default_duration *float64 `json:"default_duration"`
	Default_tz       *string `json:"default_tz" binding:"required"`
}

// получение настроек
type Setting struct {
	DefaultTz       string  `json:"default_tz" binding:"required"`
	DefaultDuration float64 `json:"default_duration" binding:"required"`
}