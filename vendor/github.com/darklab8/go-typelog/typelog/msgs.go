package typelog

import (
	"fmt"
	"log/slog"
	"os"
)

func (l *Logger) Debug(msg string, opts ...LogType) {
	if IsMsgEnabled(l.level_log, LEVEL_DEBUG) {
		args := append([]SlogAttr{}, newSlogArgs(opts...)...)
		if l.enable_file_showing {
			args = append(args, logGroupFiles())
		}
		l.logger.Debug(msg, args...)
	}

}

func (l *Logger) Info(msg string, opts ...LogType) {
	if IsMsgEnabled(l.level_log, LEVEL_INFO) {
		args := append([]SlogAttr{}, newSlogArgs(opts...)...)
		if l.enable_file_showing {
			args = append(args, logGroupFiles())
		}
		l.logger.Info(msg, args...)
	}

}

// Warn is for just potentially bad behavior to be aware of
func (l *Logger) Warn(msg string, opts ...LogType) {
	if IsMsgEnabled(l.level_log, LEVEL_WARN) {
		args := append([]SlogAttr{}, newSlogArgs(opts...)...)
		if l.enable_file_showing {
			args = append(args, logGroupFiles())
		}
		l.logger.Warn(msg, args...)
	}

}

// Error is bad but program can recover from it
func (l *Logger) Error(msg string, opts ...LogType) {
	if IsMsgEnabled(l.level_log, LEVEL_ERROR) {
		args := append([]SlogAttr{}, newSlogArgs(opts...)...)
		if l.enable_file_showing {
			args = append(args, logGroupFiles())
		}
		l.logger.Error(msg, args...)
	}

}

// Fatal when encountered, Program is not allowed to run further with fatal
func (l *Logger) Fatal(msg string, opts ...LogType) {
	args := append([]SlogAttr{}, newSlogArgs(opts...)...)
	if l.enable_file_showing {
		args = append(args, logGroupFiles())
	}
	l.logger.Error(msg, args...)

	os.Exit(1)
}

func (l *Logger) Panic(msg string, opts ...LogType) {
	args := append([]SlogAttr{}, newSlogArgs(opts...)...)
	if l.enable_file_showing {
		args = append(args, logGroupFiles())
	}
	l.logger.Error(msg, args...)
	panic(msg)
}

func (l *Logger) CheckDebug(err error, msg string, opts ...LogType) bool {
	if err == nil {
		return false
	}
	if !IsMsgEnabled(l.level_log, LEVEL_DEBUG) {
		return true
	}
	args := append([]SlogAttr{}, newSlogArgs(opts...)...)
	args = append(args, slog.String("error", fmt.Sprintf("%v", err)))
	l.logger.Debug(msg, args...)
	return true
}

func (l *Logger) CheckWarn(err error, msg string, opts ...LogType) bool {
	if err == nil {
		return false
	}
	if !IsMsgEnabled(l.level_log, LEVEL_WARN) {
		return true
	}
	args := append([]SlogAttr{}, newSlogArgs(opts...)...)
	args = append(args, slog.String("error", fmt.Sprintf("%v", err)))
	l.logger.Warn(msg, args...)
	return true
}

func (l *Logger) CheckError(err error, msg string, opts ...LogType) bool {
	if err == nil {
		return false
	}
	if !IsMsgEnabled(l.level_log, LEVEL_ERROR) {
		return true
	}
	args := append([]SlogAttr{}, newSlogArgs(opts...)...)
	args = append(args, slog.String("error", fmt.Sprintf("%v", err)))
	l.logger.Error(msg, args...)
	return true
}

// CheckFatal has shorter error output in comparison to CheckPanic
func (l *Logger) CheckFatal(err error, msg string, opts ...LogType) {
	if err == nil {
		return
	}
	args := append([]SlogAttr{}, newSlogArgs(opts...)...)
	args = append(args, slog.String("error", fmt.Sprintf("%v", err)))
	l.logger.Error(msg, args...)
	os.Exit(1)
}

func (l *Logger) CheckPanic(err error, msg string, opts ...LogType) {
	if err == nil {
		return
	}
	args := append([]SlogAttr{}, newSlogArgs(opts...)...)
	args = append(args,
		slog.String("err_msg", fmt.Sprintf("%v", err)),
		slog.String("err_type", fmt.Sprintf("%T", err)),
	)
	l.logger.Error(msg, args...)
	panic_logger.Error(msg, args...)
	panic(panic_str.String())
}
