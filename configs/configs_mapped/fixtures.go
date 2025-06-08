package configs_mapped

import (
	"context"

	"github.com/darklab8/fl-darkstat/configs/configs_settings"
)

var parsed *MappedConfigs = nil

func TestFixtureConfigs() *MappedConfigs {
	ctx := context.Background()
	if parsed != nil {
		return parsed
	}

	game_location := configs_settings.Env.FreelancerFolder
	parsed = NewMappedConfigs().Read(ctx, game_location)

	return parsed
}
