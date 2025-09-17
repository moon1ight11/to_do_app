package settings

import "time"

type Setting struct {
	DefaultTz       time.Location
	DefaultDuration time.Duration
}
