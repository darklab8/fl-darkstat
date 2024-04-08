package typelog

type LogLevel string

func (l LogLevel) ToStr() string { return string(l) }

const (
	LEVEL_DEBUG        LogLevel = "DEBUG"
	LEVEL_INFO         LogLevel = "INFO"
	LEVEL_WARN         LogLevel = "WARN"
	LEVEL_ERROR        LogLevel = "ERROR"
	LEVEL_FATAL        LogLevel = "FATAL"
	LEVEL_DEFAULT_WARN LogLevel = ""
)

func LevelToInt(level LogLevel) int {
	switch level {
	case LEVEL_DEBUG:
		return 10
	case LEVEL_INFO:
		return 20
	case LEVEL_WARN:
		return 30
	case LEVEL_ERROR:
		return 40
	case LEVEL_FATAL:
		return 50
	case LEVEL_DEFAULT_WARN:
		return 30
	}
	panic("not supported log level")
}

func IsMsgEnabled(current_level, msg_level LogLevel) bool {
	return LevelToInt(current_level) <= LevelToInt(msg_level)
}
