package configs_mapped

import "github.com/darklab8/fl-configs/configs/configs_mapped/configs_fixtures"

var parsed *MappedConfigs = nil

func TestFixtureConfigs() *MappedConfigs {
	if parsed != nil {
		return parsed
	}

	game_location := configs_fixtures.FixtureGameLocation()
	parsed = NewMappedConfigs().Read(game_location)

	return parsed
}
