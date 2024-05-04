package fronttypes

import (
	"github.com/darklab8/fl-configs/configs/configs_export"
	"github.com/darklab8/fl-configs/configs/discovery/techcompat"
)

type DiscoveryIDs struct {
	Show   bool
	Ids    []configs_export.Tractor
	Config *techcompat.Config
}
