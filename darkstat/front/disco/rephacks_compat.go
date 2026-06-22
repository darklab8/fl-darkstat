package disco

import (
	"encoding/json"
	"strings"

	"github.com/darklab8/fl-darkstat/configs/cfg"
	"github.com/darklab8/fl-darkstat/configs/discovery/playercntl_rephacks"
	"github.com/darklab8/fl-darkstat/darkstat/configs_export"
	"github.com/darklab8/fl-darkstat/darkstat/front/types"
)

// func (r RepType) ToStr() string {
// 	switch r {
// 	case MODE_REP_LESSTHAN:
// 		return "<= Maximum possible rep (not greater than)"
// 	case MODE_REP_GREATERTHAN:
// 		return ">= Minimum possible rep (not less than)"
// 	case MODE_REP_NO_CHANGE:
// 		return "? Not identified, MODE_REP_NO_CHANGE"
// 	case MODE_REP_STATIC:
// 		return "= Fixed value rep (forced static rep)"
// 	}
// 	return "undefined"
// }

func GetRephackCompat(shared *types.SharedData, faction_nicks []string, Nickname cfg.TractorID) (rep float64) {
	rep = 1.0
	tractor := shared.TractorsByID[Nickname]
	if tractor == nil || Nickname == "" {
		return 1
	}

	for _, faction_nick := range faction_nicks {

		rephack, ok := tractor.DiscoveryIDRephacks.Rephacks[cfg.FactionNick(faction_nick)]
		if !ok {
			continue
			// logus.Log.Panic(fmt.Sprintf("haven't found rephacks info for id=%s", Nickname))
		}

		switch rephack.RepType {
		case playercntl_rephacks.MODE_REP_LESSTHAN:
			if rephack.Reputation < rep {
				rep = rephack.Reputation
			}
		case playercntl_rephacks.MODE_REP_STATIC:
			if rephack.Reputation < rep {
				rep = rephack.Reputation
			}
		}

	}

	DOCKING_POSIBILITY_THRESHOLD := -0.54

	if rep < DOCKING_POSIBILITY_THRESHOLD {
		return 0
	}

	// return normalized value from 0 to 1 to make it compatible with code of tech compat
	// rep changes from -1 to +1, so next will be normalized:
	return (rep + 1) / 2
}

func MarshalRephacksCompat(shared *types.SharedData, faction_nicks []string, system_nicks []string) (result string) {
	var compat_by_id map[string]float64 = make(map[string]float64)

	compat_by_id[""] = GetRephackCompat(shared, faction_nicks, "")

	for _, id := range shared.Ids {

		is_zoner_whale_forbidden_route := false
		if id.Nickname == "dsy_license_gd_z_grp" {
			for _, system_nick := range system_nicks {
				if _, ok := configs_export.ZonerForbiddenSystems[system_nick]; ok {
					is_zoner_whale_forbidden_route = true
				}
			}
		}
		if is_zoner_whale_forbidden_route {
			continue
		}

		compat := GetRephackCompat(shared, faction_nicks, id.Nickname)

		// data size saving
		if compat <= 0.1 {
			continue
		}

		compat_by_id[string(id.ShortNickname)] = compat
	}

	bytes, _ := json.Marshal(compat_by_id)
	return strings.ReplaceAll(string(bytes), "\"", "'")
}
