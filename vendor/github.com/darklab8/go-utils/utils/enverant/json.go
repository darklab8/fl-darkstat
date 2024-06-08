package enverant

import (
	"bufio"
	"bytes"
	"encoding/json"
	"os"
	"regexp"

	"github.com/darklab8/go-typelog/typelog"
	"github.com/darklab8/go-utils/utils/regexy"
	"github.com/darklab8/go-utils/utils/utils_logus"
)

var regexConfiglines *regexp.Regexp

func init() {
	regexy.InitRegexExpression(&regexConfiglines, `^(.*)(?:// .*)$`)
}

func ReadJson(path string) map[string]interface{} {
	env_map := make(map[string]interface{})

	file, err := os.Open(path)
	if utils_logus.Log.CheckWarn(err, "not found env file at path", typelog.String("path", path)) {
		return env_map
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	var cleaned_json bytes.Buffer

	for scanner.Scan() {
		input_line := scanner.Text()
		match := regexConfiglines.FindStringSubmatch(input_line)
		if len(match) > 0 {
			input_line = match[1]
		}
		cleaned_json.WriteString(input_line)
	}

	err = json.Unmarshal(cleaned_json.Bytes(), &env_map)
	utils_logus.Log.CheckPanic(err, "failed to parse json with env vars")
	return env_map
}
