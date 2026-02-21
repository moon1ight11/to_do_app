package settings

type Setting struct {
	DefaultTz       string  `json:"default_tz" binding:"required"`
	DefaultDuration float64 `json:"default_duration" binding:"required"`
}
