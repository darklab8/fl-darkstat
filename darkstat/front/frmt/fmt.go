package frmt

/*
Allowed to be imported by anything
*/

import (
	"sort"

	"github.com/darklab8/fl-configs/configs/cfgtype"
	"github.com/darklab8/fl-configs/configs/configs_export"
)

func SortedBases(bases_map map[cfgtype.BaseUniNick]*configs_export.GoodAtBase) []*configs_export.GoodAtBase {
	var bases []*configs_export.GoodAtBase = make([]*configs_export.GoodAtBase, 0, 10)

	for _, base := range bases_map {
		bases = append(bases, base)
	}

	sort.Slice(bases, func(i, j int) bool {
		if bases[i].BaseName != "" && bases[j].BaseName == "" {
			return true
		}
		return bases[i].BaseName < bases[j].BaseName
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
