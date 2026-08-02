package inireader

import (
	"fmt"
	"regexp"

	"github.com/darklab8/fl-darkstat/configs/configs_settings/logus"

	"github.com/darklab8/go-utils/utils/utils_logus"
	"github.com/darklab8/go-utils/utils/utils_os"
)

func InitRegexExpression(regex **regexp.Regexp, expression string) {
	var err error

	if regex == nil {
		panic(fmt.Sprintln("expression has wrong pointer to memory, expression=", expression))
	}

	*regex, err = regexp.Compile(expression)
	logus.Log.CheckPanic(err, "failed to parse numberParser in ", utils_logus.FilePath(utils_os.GetCurrentFile()))
}
