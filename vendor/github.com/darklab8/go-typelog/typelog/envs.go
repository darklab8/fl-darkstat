package typelog

import (
	"os"
	"strings"
)

const (
	TOOL_NAME = "typelog"
)

var EnvTurnJSON bool = os.Getenv(strings.ToUpper(TOOL_NAME)+"_LOG_JSON") == "true"

var EnvTurnFileShowing bool = os.Getenv(strings.ToUpper(TOOL_NAME)+"_LOG_FILE_SHOWING") == "true"
