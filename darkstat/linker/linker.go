package linker

/*
Links data from exported fl-configs
into stuff rendered by fl-darkstat
*/

import (
	"sort"

	"github.com/darklab8/fl-configs/configs/configs_export"
	"github.com/darklab8/fl-configs/configs/configs_mapped"
	"github.com/darklab8/fl-darkstat/darkstat/builder"
	"github.com/darklab8/fl-darkstat/darkstat/front"
	"github.com/darklab8/fl-darkstat/darkstat/front/urls"
	"github.com/darklab8/fl-darkstat/darkstat/settings"
	"github.com/darklab8/fl-darkstat/darkstat/settings/logus"
	"github.com/darklab8/go-utils/goutils/utils/utils_logus"
	"github.com/darklab8/go-utils/goutils/utils/utils_types"
)

type Linker struct {
	configs *configs_export.Exporter
}

type LinkOption func(l *Linker)

func NewLinker(opts ...LinkOption) *Linker {
	l := &Linker{}
	for _, opt := range opts {
		opt(l)
	}

	if l.configs == nil {
		configs := configs_mapped.NewMappedConfigs()
		logus.Log.Debug("scanning freelancer folder", utils_logus.FilePath(settings.FreelancerFolder))
		configs.Read(settings.FreelancerFolder)
		l.configs = configs_export.NewExporter(configs)
	}

	return l
}

func (l *Linker) Link() *builder.Builder {
	data := l.configs.Export()

	sort.Slice(data.Bases, func(i, j int) bool {
		if data.Bases[i].Name != "" && data.Bases[j].Name == "" {
			return true
		}
		return data.Bases[i].Name < data.Bases[j].Name
	})

	for _, base := range data.Bases {
		sort.Slice(base.MarketGoods, func(i, j int) bool {
			if base.MarketGoods[i].Name != "" && base.MarketGoods[j].Name == "" {
				return true
			}
			return base.MarketGoods[i].Name < base.MarketGoods[j].Name
		})
	}

	sort.Slice(data.Factions, func(i, j int) bool {
		if data.Factions[i].Name != "" && data.Factions[j].Name == "" {
			return true
		}
		return data.Factions[i].Name < data.Factions[j].Name
	})

	for _, faction := range data.Factions {
		sort.Slice(faction.Reputations, func(i, j int) bool {
			if faction.Reputations[i].Name != "" && faction.Reputations[j].Name == "" {
				return true
			}
			return faction.Reputations[i].Name < faction.Reputations[j].Name
		})
	}

	build := builder.NewBuilder()
	build.RegComps(
		builder.NewComponent(
			urls.Index,
			front.Index(),
		),
		builder.NewComponent(
			urls.Bases,
			front.BasesT(data.Bases),
		),
		builder.NewComponent(
			urls.Factions,
			front.FactionsT(data.Factions),
		),
	)

	var infocard_per_good_nickname map[string]configs_export.Infocard = make(map[string]configs_export.Infocard)

	for _, base := range data.Bases {
		// fmt.Println("market_goods, len=", len(base.MarketGoods), " nickname=", base.Nickname, base.Name)
		build.RegComps(
			builder.NewComponent(
				utils_types.FilePath(front.BaseInfocardUrl(base)),
				front.Infocard(base.Infocard),
			),

			builder.NewComponent(
				utils_types.FilePath(front.BaseMarketGoodUrl(base)),
				front.BaseMarketGoods(base.MarketGoods),
			),
		)

		for _, good := range base.MarketGoods {
			infocard_per_good_nickname[good.Nickname] = good.Infocard
		}
	}

	for good_nickname, infocard := range infocard_per_good_nickname {
		build.RegComps(
			builder.NewComponent(
				utils_types.FilePath(front.MarketGoodInfocardUrl(good_nickname)),
				front.Infocard(infocard),
			),
		)
	}

	for _, faction := range data.Factions {
		build.RegComps(
			builder.NewComponent(
				utils_types.FilePath(front.FactionInfocardUrl(faction.Nickname)),
				front.Infocard(faction.Infocard),
			),

			builder.NewComponent(
				utils_types.FilePath(front.FactionRepUrl(faction)),
				front.FactionReps(faction.Reputations),
			),
		)
	}

	return build
}
