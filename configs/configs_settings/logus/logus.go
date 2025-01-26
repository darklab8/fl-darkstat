package logus

import (
	_ "github.com/darklab8/fl-darkstat/configs/configs_settings"
	"github.com/darklab8/go-typelog/typelog"
)

var Log *typelog.Logger = typelog.NewLogger("configs")
