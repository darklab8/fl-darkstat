package regexy

import (
	"fmt"
	"regexp"

	"github.com/darklab8/go-utils/utils/utils_types"
)

func InitRegexExpression(regex **regexp.Regexp, expression utils_types.RegExp) {
	var err error

	*regex, err = regexp.Compile(string(expression))

	if err != nil {
		panic(fmt.Sprintln(err, "failed to init regex=", expression))
	}
}
