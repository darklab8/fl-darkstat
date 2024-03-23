package logus

import (
	"github.com/darklab8/go-typelog/typelog"
)

var Log *typelog.Logger = typelog.NewLogger("darkstat",
	typelog.WithLogLevel(typelog.LEVEL_INFO),
)
