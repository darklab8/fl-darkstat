package inireader

import (
	"regexp"

	"github.com/darklab8/fl-darkstat/configs/configs_settings/logus"

	"github.com/darklab8/go-utils/utils/utils_logus"
	"github.com/darklab8/go-utils/utils/utils_os"
)

func InitRegexExpression(regex **regexp.Regexp, expression string) {
	var err error

	*regex, err = regexp.Compile(expression)
	logus.Log.CheckPanic(err, "failed to parse numberParser in ", utils_logus.FilePath(utils_os.GetCurrentFile()))
}
