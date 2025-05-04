// package logus - logging confugiration for darkcore. Useful for global usage or getting childs into modules
package logus

import (
	_ "github.com/darklab8/fl-darkstat/darkcore/settings"
	"github.com/darklab8/go-typelog/typelog"
)

var Log *typelog.Logger = typelog.NewLogger("darkcore",
	typelog.WithLogLevel(typelog.LEVEL_INFO),
)
