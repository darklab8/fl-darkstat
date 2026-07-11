package frmt

/*
Allowed to be imported by anything
*/

import (
	"sort"

	"github.com/darklab8/fl-darkstat/configs/cfg"
	"github.com/darklab8/fl-darkstat/darkstat/configs_export"
)

func CompareBool(a bool, b bool) int {
	if !a && b {
		return -1
	}
	if a && !b {
		return 1
	}
	return 0
}

func SortedBases(bases_map map[cfg.BaseUniNick]*configs_export.MarketGood) []*configs_export.MarketGood {
	var bases []*configs_export.MarketGood = make([]*configs_export.MarketGood, 0, 10)

	for _, base := range bases_map {
		bases = append(bases, base)
	}

	sort.Slice(bases, func(i, j int) bool {
		if bases[i].BaseSells != bases[j].BaseSells {
			return bases[i].BaseSells // true sorts before false
		}
		return bases[i].PriceBaseSellsFor < bases[j].PriceBaseSellsFor
	})

	return bases
}

func FormatBaseSells(value bool) string {
	if value {
		return "sells"
	} else {
		return "buysonly"
	}
}

func FormatBoolAsYesNo(value bool) string {
	if value {
		return "yes"
	} else {
		return "no"
	}
}

func FmtPtrStringOrQuestionMark(value *string) string {
	if value == nil {
		return "?"
	}
	return *value
}
