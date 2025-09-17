package minecontrol

import (
	"github.com/darklab8/fl-darkstat/configs/cfg"
	"github.com/darklab8/fl-darkstat/configs/configs_mapped/parserutils/iniload"
	"github.com/darklab8/fl-darkstat/configs/configs_mapped/parserutils/semantic"
)

type PlayerBonus struct {
	semantic.Model
	IDNickname  *semantic.String
	OreNickname *semantic.String
	Bonus       *semantic.Float
	Index       int
}

type Config struct {
	PlayerBonusByIDNickname  map[cfg.TractorID][]*PlayerBonus
	PlayerBonusByOreNickname map[string][]*PlayerBonus
}

func Read(input_file *iniload.IniLoader) *Config {
	conf := &Config{
		PlayerBonusByIDNickname:  map[cfg.TractorID][]*PlayerBonus{},
		PlayerBonusByOreNickname: map[string][]*PlayerBonus{},
	}

	if resources, ok := input_file.SectionMap["[playerbonus]"]; ok {
		resource := resources[0]

		for index, param := range resource.ParamMap["pb"] {
			player_bonus := &PlayerBonus{
				IDNickname:  semantic.NewString(resource, param.Key, semantic.OptsS(semantic.Order(0), semantic.Index(index))),
				OreNickname: semantic.NewString(resource, param.Key, semantic.OptsS(semantic.Order(1), semantic.Index(index))),
				Bonus:       semantic.NewFloat(resource, param.Key, semantic.Precision(2), semantic.OptsF(semantic.Order(2), semantic.Index(index))),
				Index:       index,
			}
			player_bonus.Map(resource)
			conf.PlayerBonusByIDNickname[cfg.TractorID(player_bonus.IDNickname.Get())] = append(conf.PlayerBonusByIDNickname[cfg.TractorID(player_bonus.IDNickname.Get())], player_bonus)
			conf.PlayerBonusByOreNickname[player_bonus.OreNickname.Get()] = append(conf.PlayerBonusByOreNickname[player_bonus.OreNickname.Get()], player_bonus)
		}

	}
	return conf
}
