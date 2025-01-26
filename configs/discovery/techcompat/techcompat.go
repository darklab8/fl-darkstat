package techcompat

import (
	"github.com/darklab8/fl-darkstat/configs/cfg"
	"github.com/darklab8/fl-darkstat/configs/configs_mapped/parserutils/filefind/file"
	"github.com/darklab8/fl-darkstat/configs/configs_mapped/parserutils/iniload"
	"github.com/darklab8/fl-darkstat/configs/configs_mapped/parserutils/semantic"
)

type General struct {
	semantic.Model
	UnlistedTech  *semantic.Float
	DefaultMult   *semantic.Float
	NoControlItem *semantic.Float
}

type TechCompatibility struct {
	semantic.Model
	Nickname   *semantic.String
	Percentage *semantic.Float
}

type Faction struct {
	semantic.Model
	DefaultUnlisted *semantic.Float
	ID              *semantic.String
	TechCompats     []*TechCompatibility
}

type TechGroup struct {
	semantic.Model
	Name    *semantic.String
	Default *semantic.Float
	Items   []*semantic.String
}

type Config struct {
	*iniload.IniLoader
	General         *General
	Factions        []*Faction
	FactionByID     map[cfg.TractorID]*Faction
	TechGroups      []*TechGroup
	TechGroupByName map[string]*TechGroup

	// string ItemNickname
	CompatByItem map[string]*ItemCompat
}

type ItemCompat struct {
	Default    *float64
	TechCell   string
	CompatByID map[cfg.TractorID]float64
}

func (conf *Config) GetCompatibilty(item_nickname string, id_nickname cfg.TractorID) float64 {

	if id_nickname == "" {
		// ; If the ship does not have a control item (in discovery this is the ID) then this
		// ; multiplier is used.
		return conf.General.NoControlItem.Get()
	}

	faction, is_found_faction := conf.FactionByID[id_nickname]

	item, found_item := conf.CompatByItem[item_nickname]
	if !found_item {
		if is_found_faction {
			if default_unlisted, ok := faction.DefaultUnlisted.GetValue(); ok {
				return default_unlisted
			}
		}
		// ; Any items not in a [tech] section use this multiplier.
		// unlisted_tech = smth
		return conf.General.UnlistedTech.Get()
	}

	item_faction_compat, found_faction := item.CompatByID[id_nickname]

	if !found_faction {
		// ; Anything in listed in a [tech] section but not an explicitly defined in the faction
		// ; section combination uses this as the default multipier.

		if item.Default != nil {
			return *item.Default
		} else {
			return conf.General.DefaultMult.Get()
		}
	}

	return item_faction_compat
}

func Read(input_file *iniload.IniLoader) *Config {
	conf := &Config{
		IniLoader:       input_file,
		FactionByID:     make(map[cfg.TractorID]*Faction),
		TechGroupByName: make(map[string]*TechGroup),
		CompatByItem:    make(map[string]*ItemCompat),
	}

	if resources, ok := input_file.SectionMap["[general]"]; ok {
		general_info := resources[0]

		conf.General = &General{
			UnlistedTech:  semantic.NewFloat(general_info, cfg.Key("unlisted_tech"), semantic.Precision(2)),
			DefaultMult:   semantic.NewFloat(general_info, cfg.Key("default_mult"), semantic.Precision(2)),
			NoControlItem: semantic.NewFloat(general_info, cfg.Key("no_control_item"), semantic.Precision(2)),
		}
		conf.General.Map(general_info)

	}

	for _, faction_info := range input_file.SectionMap["[faction]"] {

		faction_nicknames := faction_info.ParamMap[cfg.Key("item")][0]
		for faction_order, _ := range faction_nicknames.Values {
			faction := &Faction{
				DefaultUnlisted: semantic.NewFloat(faction_info, cfg.Key("default_unlisted"), semantic.Precision(2)),
			}
			faction.Map(faction_info)
			faction.ID = semantic.NewString(faction_info, cfg.Key("item"), semantic.OptsS(semantic.Order(faction_order)))

			for index, _ := range faction_info.ParamMap[cfg.Key("tech")] {
				compat := &TechCompatibility{}
				compat.Map(faction_info)
				compat.Nickname = semantic.NewString(faction_info, cfg.Key("tech"), semantic.OptsS(semantic.Index(index), semantic.Order(0)))
				compat.Percentage = semantic.NewFloat(faction_info, cfg.Key("tech"), semantic.Precision(2), semantic.OptsF(semantic.Index(index), semantic.Order(1)))
				faction.TechCompats = append(faction.TechCompats, compat)
			}

			conf.Factions = append(conf.Factions, faction)
			conf.FactionByID[cfg.TractorID(faction.ID.Get())] = faction
		}
	}

	for _, techgroup_info := range input_file.SectionMap["[tech]"] {
		techgroup := &TechGroup{}
		techgroup.Map(techgroup_info)

		techgroup.Name = semantic.NewString(techgroup_info, cfg.Key("name"))
		techgroup.Default = semantic.NewFloat(techgroup_info, cfg.Key("default"), semantic.Precision(2))

		for index, _ := range techgroup_info.ParamMap[cfg.Key("item")] {
			techgroup.Items = append(techgroup.Items, semantic.NewString(techgroup_info, cfg.Key("item"), semantic.OptsS(semantic.Index(index))))
		}

		conf.TechGroups = append(conf.TechGroups, techgroup)
		conf.TechGroupByName[techgroup.Name.Get()] = techgroup

		for _, item := range techgroup.Items {
			item_nickname := item.Get()
			compat, found_compat := conf.CompatByItem[item_nickname]

			if !found_compat {
				compat = &ItemCompat{CompatByID: make(map[cfg.TractorID]float64)}
				conf.CompatByItem[item_nickname] = compat

				if value, ok := techgroup.Default.GetValue(); ok {
					compat.Default = &value
				}
			}

			compat.TechCell = techgroup.Name.Get()

			for _, faction := range conf.Factions {
				for _, faction_compat := range faction.TechCompats {
					if compat.TechCell != faction_compat.Nickname.Get() {
						continue
					}

					id_nickname := cfg.TractorID(faction.ID.Get())
					compat.CompatByID[id_nickname] = faction_compat.Percentage.Get()

				}
			}
		}
	}

	return conf
}

func (frelconfig *Config) Write() *file.File {
	return &file.File{}
}
