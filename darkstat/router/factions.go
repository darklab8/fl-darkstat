package router

import (
	"sort"

	"github.com/darklab8/fl-darkstat/darkcore/builder"
	"github.com/darklab8/fl-darkstat/darkstat/configs_export"
	"github.com/darklab8/fl-darkstat/darkstat/front"
	"github.com/darklab8/fl-darkstat/darkstat/front/tab"
	"github.com/darklab8/fl-darkstat/darkstat/front/types"
	"github.com/darklab8/fl-darkstat/darkstat/front/urls"
	"github.com/darklab8/go-utils/utils/timeit"
	"github.com/darklab8/go-utils/utils/utils_types"
)

func (l *Router) LinkFactions(
	build *builder.Builder,
	data *configs_export.Exporter,
	shared *types.SharedData,
) {
	sort.Slice(data.Factions, func(i, j int) bool {
		if data.Factions[i].Name != "" && data.Factions[j].Name == "" {
			return true
		}
		return data.Factions[i].Name < data.Factions[j].Name
	})

	for fac_index, faction := range data.Factions {
		var reps []configs_export.Reputation = make([]configs_export.Reputation, 0, len(faction.Reputations))
		for _, rep := range faction.Reputations {
			if rep.Name != "" {
				reps = append(reps, rep)
			}
		}
		sort.Slice(reps, func(i, j int) bool {
			return reps[i].Rep > reps[j].Rep
		})
		data.Factions[fac_index].Reputations = reps
	}

	var useful_factions []configs_export.Faction = configs_export.FilterToUsefulFactions(data.Factions)
	var useful_bribes []configs_export.Faction = configs_export.FilterToUsefulBribes(data.Factions)

	build.RegComps(
		builder.NewComponent(
			urls.Factions,
			front.FactionsT(useful_factions, front.FactionShowBases, tab.ShowEmpty(false), shared),
		),
		builder.NewComponent(
			tab.AllItemsUrl(urls.Factions),
			front.FactionsT(data.Factions, front.FactionShowBases, tab.ShowEmpty(true), shared),
		),
		builder.NewComponent(
			urls.Bribes,
			front.FactionsT(useful_bribes, front.FactionShowBribes, tab.ShowEmpty(false), shared),
		),
		builder.NewComponent(
			tab.AllItemsUrl(urls.Bribes),
			front.FactionsT(data.Factions, front.FactionShowBribes, tab.ShowEmpty(true), shared),
		),
	)

	timeit.NewTimerMF("linking faction stuff", func() {
		for _, faction := range data.Factions {
			build.RegComps(
				builder.NewComponent(
					utils_types.FilePath(front.FactionRepUrl(faction, front.FactionShowBases)),
					front.FactionReps(faction, faction.Reputations),
				),
				builder.NewComponent(
					utils_types.FilePath(front.FactionRepUrl(faction, front.FactionShowBribes)),
					front.RephackBottom(faction, faction.Bribes),
				),
			)
		}
	})
}
