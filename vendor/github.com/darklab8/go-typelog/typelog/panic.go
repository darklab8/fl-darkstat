package typelog

import (
	"bytes"
	"log/slog"
)

/*
Global logger for writing formatted msg into string we could input to panic.
The point of this code that we get string with all formatted slog.Attrs for panic.
*/

var (
	panic_logger *slog.Logger
	panic_str    *bytes.Buffer = bytes.NewBuffer([]byte{})
)

func init() {
	panic_logger = slog.New(slog.NewTextHandler(panic_str, &slog.HandlerOptions{Level: slog.LevelDebug}))
}
