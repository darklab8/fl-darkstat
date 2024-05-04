package techcompat

import (
	"github.com/darklab8/fl-configs/configs/configs_mapped/parserutils/filefind/file"
	"github.com/darklab8/fl-configs/configs/configs_mapped/parserutils/iniload"
	"github.com/darklab8/fl-configs/configs/configs_mapped/parserutils/semantic"
	"github.com/darklab8/fl-configs/configs/conftypes"
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
	ID          *semantic.String
	TechCompats []*TechCompatibility
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
	FactionByID     map[string]*Faction
	TechGroups      []*TechGroup
	TechGroupByName map[string]*TechGroup

	// string ItemNickname
	CompatByItem map[string]*ItemCompat
}

type ItemCompat struct {
	Default    *float64
	TechCell   string
	CompatByID map[conftypes.TractorID]float64
}

func (conf *Config) GetCompatibilty(item_nickname string, id_nickname conftypes.TractorID) float64 {

	if id_nickname == "" {
		// ; If the ship does not have a control item (in discovery this is the ID) then this
		// ; multiplier is used.
		return conf.General.NoControlItem.Get()
	}

	item, found_item := conf.CompatByItem[item_nickname]
	if !found_item {
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
		FactionByID:     make(map[string]*Faction),
		TechGroupByName: make(map[string]*TechGroup),
		CompatByItem:    make(map[string]*ItemCompat),
	}

	if resources, ok := input_file.SectionMap["[general]"]; ok {
		general_info := resources[0]

		conf.General = &General{
			UnlistedTech:  semantic.NewFloat(general_info, "unlisted_tech", semantic.Precision(2)),
			DefaultMult:   semantic.NewFloat(general_info, "default_mult", semantic.Precision(2)),
			NoControlItem: semantic.NewFloat(general_info, "no_control_item", semantic.Precision(2)),
		}
		conf.General.Map(general_info)

	}

	for _, faction_info := range input_file.SectionMap["[faction]"] {
		faction := &Faction{}
		faction.Map(faction_info)

		faction.ID = semantic.NewString(faction_info, "item")

		for index, _ := range faction_info.ParamMap["tech"] {
			compat := &TechCompatibility{}
			compat.Map(faction_info)
			compat.Nickname = semantic.NewString(faction_info, "tech", semantic.OptsS(semantic.Index(index), semantic.Order(0)))
			compat.Percentage = semantic.NewFloat(faction_info, "tech", semantic.Precision(2), semantic.Index(index), semantic.Order(1))
			faction.TechCompats = append(faction.TechCompats, compat)
		}

		conf.Factions = append(conf.Factions, faction)
		conf.FactionByID[faction.ID.Get()] = faction
	}

	for _, techgroup_info := range input_file.SectionMap["[tech]"] {
		techgroup := &TechGroup{}
		techgroup.Map(techgroup_info)

		techgroup.Name = semantic.NewString(techgroup_info, "name")
		techgroup.Default = semantic.NewFloat(techgroup_info, "default", semantic.Precision(2))

		for index, _ := range techgroup_info.ParamMap["item"] {
			techgroup.Items = append(techgroup.Items, semantic.NewString(techgroup_info, "item", semantic.OptsS(semantic.Index(index))))
		}

		conf.TechGroups = append(conf.TechGroups, techgroup)
		conf.TechGroupByName[techgroup.Name.Get()] = techgroup

		for _, item := range techgroup.Items {
			item_nickname := item.Get()
			compat, found_compat := conf.CompatByItem[item_nickname]

			if !found_compat {
				compat = &ItemCompat{CompatByID: make(map[conftypes.TractorID]float64)}
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

					id_nickname := conftypes.TractorID(faction.ID.Get())
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
