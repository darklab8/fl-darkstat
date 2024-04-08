package utils

import (
	"regexp"

	"github.com/darklab8/go-utils/goutils/utils/utils_logus"
	"github.com/darklab8/go-utils/goutils/utils/utils_types"
)

func InitRegexExpression(regex **regexp.Regexp, expression utils_types.RegExp) {
	var err error

	*regex, err = regexp.Compile(string(expression))
	utils_logus.Log.CheckFatal(err, "failed to init regex",
		utils_logus.Regex(expression), utils_logus.FilePath(GetCurrentFile()))
}
