package configs_mapped

import (
	"github.com/darklab8/fl-darkstat/configs/configs_settings"
)

var parsed *MappedConfigs = nil

func TestFixtureConfigs() *MappedConfigs {
	if parsed != nil {
		return parsed
	}

	game_location := configs_settings.Env.FreelancerFolder
	parsed = NewMappedConfigs().Read(game_location)

	return parsed
}
