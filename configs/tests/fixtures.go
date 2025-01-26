package tests

import (
	"github.com/darklab8/fl-darkstat/configs/configs_mapped/parserutils/filefind"
	"github.com/darklab8/fl-darkstat/configs/configs_settings"
)

var cached *filefind.Filesystem

func FixtureFileFind() *filefind.Filesystem {
	if cached != nil {
		return cached
	}
	cached = filefind.FindConfigs(configs_settings.Env.FreelancerFolder)
	return cached
}
