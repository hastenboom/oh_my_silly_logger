package core

import "time"

type Level uint8

const (
	DEBUG Level = iota
	INFO
	WARN
	ERROR
	FATAL
)

//expansion of time.Duration enum
const (
	DAY   time.Duration = 24 * time.Hour
	MONTH time.Duration = 30 * DAY
)

func getLevelString(level Level) string {
	switch level {
	case DEBUG:
		return "DEBUG"
	case INFO:
		return "INFO"
	case WARN:
		return "WARN"
	case ERROR:
		return "ERROR"
	case FATAL:
		return "FATAL"
	default:
		return "UNKNOWN"
	}
}
