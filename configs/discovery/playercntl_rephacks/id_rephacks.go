package playercntl_rephacks

import (
	"github.com/darklab8/fl-darkstat/configs/cfg"
	"github.com/darklab8/fl-darkstat/configs/configs_mapped/parserutils/iniload"
	"github.com/darklab8/fl-darkstat/configs/configs_mapped/parserutils/inireader"
	"github.com/darklab8/fl-darkstat/configs/configs_mapped/parserutils/semantic"
)

type RepType int

func (r RepType) ToStr() string {
	switch r {
	case MODE_REP_LESSTHAN:
		return "<= Maximum possible rep (not greater than)"
	case MODE_REP_GREATERTHAN:
		return ">= Minimum possible rep (not less than)"
	case MODE_REP_NO_CHANGE:
		return "? Not identified, MODE_REP_NO_CHANGE"
	case MODE_REP_STATIC:
		return "= Fixed value rep (forced static rep)"
	}
	return "undefined"
}

const (
	// COPIED FROM https://github.com/Aingar/FLHook/blob/ca76a9fbfb74c5c5d609bd5042adca45a5ee866c/Plugins/Public/playercntl_plugin/RepFixer.cpp#L40

	// The adjustment mode. If the player's reputation for scRepGroup
	// is greater than fRep then make the reputation equal to fRep
	MODE_REP_LESSTHAN RepType = 0

	// The adjustment mode. If the player's reputation for scRepGroup
	// is less than fRep then make the reputation equal to fRep
	MODE_REP_GREATERTHAN RepType = 1

	// Don't change anything/ignore this reputation group.
	MODE_REP_NO_CHANGE RepType = 2

	// Fix the rep group to this level.
	MODE_REP_STATIC RepType = 3
)

type Faction struct {
	semantic.Model
	Nickname cfg.FactionNick

	Rep *semantic.Float

	// unknown: supposedly 0 means minimum or exact. 3 means maximum or exact.
	RepType *semantic.Int
}

func (f Faction) GetRepType() RepType {
	return RepType(f.RepType.Get())
}

type Rephack struct {
	semantic.Model
	ID       *semantic.String
	Inherits *semantic.String
	Reps     map[cfg.FactionNick]Faction
}

type Config struct {
	DefaultReps map[cfg.FactionNick]Faction

	RephacksByID map[cfg.TractorID]Rephack
}

func Read(input_file *iniload.IniLoader) *Config {
	conf := &Config{
		DefaultReps:  make(map[cfg.FactionNick]Faction),
		RephacksByID: make(map[cfg.TractorID]Rephack),
	}

	if resources, ok := input_file.SectionMap["[default_reps]"]; ok {

		default_reps := resources[0]

		for _, param := range default_reps.Params {
			faction := Faction{
				Nickname: cfg.FactionNick(param.Key),
				Rep:      semantic.NewFloat(default_reps, param.Key, semantic.Precision(2)),
				RepType:  semantic.NewInt(default_reps, param.Key, semantic.Order(1)),
			}
			faction.Map(default_reps)
			conf.DefaultReps[cfg.FactionNick(param.Key)] = faction
		}
	}

	for _, rephack_info := range input_file.SectionMap["[rephack]"] {
		rephack := Rephack{
			ID:       semantic.NewString(rephack_info, cfg.Key("id")),
			Inherits: semantic.NewString(rephack_info, cfg.Key("inherits")),
			Reps:     make(map[cfg.FactionNick]Faction),
		}
		rephack.Map(rephack_info)

		for _, param := range rephack_info.Params {
			if param.Key == cfg.Key("id") || param.Key == cfg.Key("inherits") || param.Key == inireader.KEY_COMMENT {
				continue
			}

			faction := Faction{
				Nickname: cfg.FactionNick(param.Key),
				Rep:      semantic.NewFloat(rephack_info, param.Key, semantic.Precision(2)),
				RepType:  semantic.NewInt(rephack_info, param.Key, semantic.Order(1)),
			}
			faction.Map(rephack_info)
			rephack.Reps[cfg.FactionNick(param.Key)] = faction
		}

		conf.RephacksByID[cfg.TractorID(rephack.ID.Get())] = rephack
	}
	return conf
}
