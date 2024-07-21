/*
Package typelog is a slog modified for extra static type safety
and extra boilerplating out of the box
*/
package typelog

import (
	"fmt"
	"io"
	"log/slog"
	"os"
	"strings"
)

type Logger struct {
	logger              *slog.Logger
	name                string
	enable_file_showing bool
	enable_json_format  bool
	log_level_slog      *slog.LevelVar
	level_log           LogLevel
	io_writer           io.Writer
}

type LoggerParam func(r *Logger)

func WithJsonFormat(state bool) LoggerParam {
	return func(logger *Logger) {
		logger.enable_json_format = state
	}
}

func WithIoWriter(writer io.Writer) LoggerParam {
	return func(logger *Logger) {
		logger.io_writer = writer
	}
}

func WithFileShowing(state bool) LoggerParam {
	return func(logger *Logger) {
		logger.enable_file_showing = state
	}
}

func WithLogLevelStr(log_level_str string) LoggerParam {
	return WithLogLevel(LogLevel(log_level_str))
}

func WithLogLevel(log_level_str LogLevel) LoggerParam {
	return func(logger *Logger) {
		var log_level LogLevel = LogLevel(log_level_str)

		var programLevel = new(slog.LevelVar) // Info by default
		switch log_level {
		case LEVEL_DEBUG:
			programLevel.Set(slog.LevelDebug)
		case LEVEL_INFO:
			programLevel.Set(slog.LevelInfo)
		case LEVEL_WARN:
			programLevel.Set(slog.LevelWarn)
		case LEVEL_ERROR:
			programLevel.Set(slog.LevelError)
		case LEVEL_DEFAULT_WARN:
			programLevel.Set(slog.LevelWarn)
		default:
			panic(fmt.Sprintf("invalid log level=%s for logger=%s", log_level_str, logger.name))
		}
		logger.log_level_slog = programLevel
		logger.level_log = log_level_str
	}
}

/*
RegisteredLoggers leaves option for end applications to access all registered loggers
and choosing their own log levels to override for them
*/
var RegisteredLoggers []*Logger

func (l *Logger) GetName() string { return l.name }

func NewLogger(
	name string,
	options ...LoggerParam,
) *Logger {

	logger := &Logger{
		name:      name,
		io_writer: os.Stdout,
	}
	RegisteredLoggers = append(RegisteredLoggers, logger)

	WithJsonFormat(Env.EnableJson)(logger)
	WithFileShowing(Env.EnableFileShowing)(logger)
	WithLogLevelStr(os.Getenv(strings.ToUpper(name) + "_LOG_LEVEL"))(logger)

	for _, opt := range options {
		opt(logger)
	}

	return logger.Initialized()
}

/*
OverrideOption for overrides by external libraries
*/
func (l *Logger) OverrideOption(options ...LoggerParam) *Logger {
	for _, opt := range options {
		opt(l)
	}

	return l.Initialized()
}

func (l *Logger) Initialized() *Logger {
	if l.enable_json_format {
		l.logger = slog.New(slog.NewJSONHandler(l.io_writer, &slog.HandlerOptions{Level: l.log_level_slog}))
	} else {
		l.logger = slog.New(slog.NewTextHandler(l.io_writer, &slog.HandlerOptions{Level: l.log_level_slog}))
	}
	return l
}

func (l *Logger) WithFields(opts ...LogType) *Logger {
	var newLogger Logger = *l
	newLogger.Initialized()
	newLogger.logger = newLogger.logger.With(newSlogArgs(opts...)...)
	return &newLogger
}
