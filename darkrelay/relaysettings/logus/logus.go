package logus

import (
	_ "github.com/darklab8/fl-darkstat/darkstat/settings"
	"github.com/darklab8/go-utils/typelog"
)

var Log *typelog.Logger = typelog.NewLogger("darkrelay",
	typelog.WithLogLevel(typelog.LEVEL_INFO),
)
