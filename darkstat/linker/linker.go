package linker

/*
Links data from exported fl-configs
into stuff rendered by fl-darkstat
*/

import (
	"fmt"
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
		return data.Bases[i].Name < data.Bases[j].Name
	})

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
			urls.Systems,
			front.Systems(),
		),
	)

	for _, base := range data.Bases {
		fmt.Println("market_goods, len=", len(base.MarketGoods), " nickname=", base.Nickname, base.Name)
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
	}

	goods := l.configs.GoodsSelEquip
	for _, good := range goods {
		build.RegComps(
			builder.NewComponent(
				utils_types.FilePath(front.MarketGoodInfocardUrl(good.Nickname)),
				front.Infocard(good.Infocard),
			),
		)
	}

	return build
}
