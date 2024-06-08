package regexy

import (
	"regexp"

	"github.com/darklab8/go-utils/utils/utils_logus"
	"github.com/darklab8/go-utils/utils/utils_os"
	"github.com/darklab8/go-utils/utils/utils_types"
)

func InitRegexExpression(regex **regexp.Regexp, expression utils_types.RegExp) {
	var err error

	*regex, err = regexp.Compile(string(expression))
	utils_logus.Log.CheckPanic(err, "failed to init regex",
		utils_logus.Regex(expression), utils_logus.FilePath(utils_os.GetCurrentFile()))
}
