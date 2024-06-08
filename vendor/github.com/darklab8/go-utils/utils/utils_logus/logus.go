package utils_logus

import (
	"fmt"
	"log/slog"

	"github.com/darklab8/go-typelog/typelog"
	"github.com/darklab8/go-utils/utils/utils_types"

	_ "github.com/darklab8/go-utils/utils/utils_settings"
)

var Log *typelog.Logger = typelog.NewLogger("utils")

func Regex(value utils_types.RegExp) typelog.LogType {
	return func(c *typelog.LogAtrs) {
		c.Append(slog.String("regexp", fmt.Sprintf("%v", value)))
	}
}

func FilePath(value utils_types.FilePath) typelog.LogType {
	return func(c *typelog.LogAtrs) {
		c.Append(slog.String("filepath", fmt.Sprintf("%v", value)))
	}
}

func Filepaths(values []utils_types.FilePath) typelog.LogType {
	return typelog.Items[utils_types.FilePath]("filepaths", values)
}
