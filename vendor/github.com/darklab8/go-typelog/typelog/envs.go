package typelog

import (
	"os"
	"strings"
)

const (
	TOOL_NAME = "typelog"
)

type TypelogEnvs struct {
	EnableJson        bool
	EnableFileShowing bool
}

var Env TypelogEnvs = TypelogEnvs{
	EnableJson:        os.Getenv(strings.ToUpper(TOOL_NAME)+"_LOG_JSON") == "true",
	EnableFileShowing: os.Getenv(strings.ToUpper(TOOL_NAME)+"_LOG_FILE_SHOWING") == "true",
}
