package enverant

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"regexp"
	"strconv"

	"github.com/darklab8/go-utils/utils/regexy"
)

var regexConfiglines *regexp.Regexp

func init() {
	regexy.InitRegexExpression(&regexConfiglines, `^(.*)(?:// .*)$`)
}

func ReadJson(path string) map[string]interface{} {
	env_map := make(map[string]interface{})

	file, err := os.Open(path)
	if err != nil {
		log.Println(err, "not found env file at path=", path)
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
	if err != nil {
		panic(fmt.Sprintln(err, "failed to parse json with env vars"))
	}

	for key, value := range env_map {

		if _, ok := os.LookupEnv(key); ok {
			continue
		}

		switch v := value.(type) {
		case bool:
			os.Setenv(key, strconv.FormatBool(v))
		case string:
			os.Setenv(key, v)
		case int:
			os.Setenv(key, fmt.Sprintf("%d", v))
		case float64:
			os.Setenv(key, fmt.Sprintf("%.0f", v))
		default:
			panic(fmt.Sprintln("enverant value in file has not supported type", key, fmt.Sprintf("%T", value)))
		}

	}

	return env_map
}
