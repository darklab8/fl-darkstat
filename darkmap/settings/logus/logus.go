package logus

import (
	_ "github.com/darklab8/fl-darkstat/darkmap/settings"
	"github.com/darklab8/go-typelog/typelog"
)

var Log *typelog.Logger = typelog.NewLogger("darkmap",
	typelog.WithLogLevel(typelog.LEVEL_INFO),
)
