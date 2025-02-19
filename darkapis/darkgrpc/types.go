package darkgrpc

import (
	"github.com/darklab8/fl-darkstat/configs/cfg"
	"github.com/darklab8/fl-darkstat/darkstat/configs_export"
)

type Nicknamable interface {
	GetNickname() string
}

type Marketable interface {
	Nicknamable
	GetBases() map[cfg.BaseUniNick]*configs_export.MarketGood
}

type TechCompatable interface {
	Nicknamable
	GetDiscoveryTechCompat() *configs_export.DiscoveryTechCompat
}
