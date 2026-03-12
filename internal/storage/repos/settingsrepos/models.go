package settingsrepos

type Setting struct {
	TimeDuration float64 `json:"time_duration"`
	UserTz       string  `json:"user_tz" binding:"required"`
}
