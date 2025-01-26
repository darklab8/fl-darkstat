package mbases_mapped

import (
	"github.com/darklab8/fl-darkstat/configs/cfg"
	"github.com/darklab8/fl-darkstat/configs/configs_mapped/parserutils/filefind/file"
	"github.com/darklab8/fl-darkstat/configs/configs_mapped/parserutils/iniload"
	"github.com/darklab8/fl-darkstat/configs/configs_mapped/parserutils/semantic"
	"github.com/darklab8/fl-darkstat/configs/configs_settings/logus"
)

const (
	FILENAME = "mbases.ini"
)

type Mroom struct {
	semantic.Model
	Nickname         *semantic.String
	CharacterDensity *semantic.Int
	Bartrender       *semantic.String
}

type MissionType struct {
	semantic.Model
	MinDifficulty *semantic.Float
	MaxDifficulty *semantic.Float
	Weight        *semantic.Int
}

type BaseFaction struct {
	semantic.Model

	MissionType *MissionType
	Faction     *semantic.String
	Weight      *semantic.Int
	Npcs        []*semantic.String
}

type Bribe struct {
	semantic.Model
	Faction *semantic.String
}
type Rumor struct {
	semantic.Model
}
type Mission struct {
	semantic.Model
}
type Know struct {
	semantic.Model
}

type NPC struct {
	semantic.Model

	Nickname    *semantic.String
	Room        *semantic.String
	Bribes      []*Bribe
	Rumors      []*Rumor
	Missions    []*Mission
	Knows       []*Know
	Affiliation *semantic.String
}

type MVendor struct {
	semantic.Model
	MinOffers *semantic.Int
	MaxOffers *semantic.Int
}

type Base struct {
	semantic.Model

	Nickname     *semantic.String
	LocalFaction *semantic.String
	Diff         *semantic.Int

	BaseFactions    []*BaseFaction
	BaseFactionsMap map[string]*BaseFaction
	NPCs            []*NPC
	Bar             *Mroom
	MVendor         *MVendor
}

type Config struct {
	semantic.ConfigModel

	File    *iniload.IniLoader
	Bases   []*Base
	BaseMap map[cfg.BaseUniNick]*Base
}

func Read(input_file *iniload.IniLoader) *Config {
	frelconfig := &Config{
		File:    input_file,
		Bases:   make([]*Base, 0, 100),
		BaseMap: make(map[cfg.BaseUniNick]*Base),
	}

	for i := 0; i < len(input_file.Sections); i++ {

		if input_file.Sections[i].Type == "[mbase]" {

			mbase_section := input_file.Sections[i]
			base := &Base{
				BaseFactionsMap: make(map[string]*BaseFaction),
			}
			base.Map(mbase_section)
			base.Nickname = semantic.NewString(mbase_section, cfg.Key("nickname"), semantic.WithLowercaseS(), semantic.WithoutSpacesS())
			base.LocalFaction = semantic.NewString(mbase_section, cfg.Key("local_faction"))
			base.Diff = semantic.NewInt(mbase_section, cfg.Key("diff"))
			frelconfig.Bases = append(frelconfig.Bases, base)
			frelconfig.BaseMap[cfg.BaseUniNick(base.Nickname.Get())] = base

			for j := i + 1; j < len(input_file.Sections) && input_file.Sections[j].Type != "[mbase]"; j++ {
				section := input_file.Sections[j]

				switch section.Type {
				case "[mvendor]":
					vendor := &MVendor{
						MinOffers: semantic.NewInt(section, cfg.Key("num_offers"), semantic.Order(0)),
						MaxOffers: semantic.NewInt(section, cfg.Key("num_offers"), semantic.Order(1)),
					}
					base.MVendor = vendor
				case "[basefaction]":
					faction := &BaseFaction{
						Faction: semantic.NewString(section, cfg.Key("faction"), semantic.WithLowercaseS(), semantic.WithoutSpacesS()),
						Weight:  semantic.NewInt(section, cfg.Key("weight")),
					}
					faction.Map(section)

					mission_type := &MissionType{
						MinDifficulty: semantic.NewFloat(section, cfg.Key("mission_type"), semantic.Precision(2), semantic.OptsF(semantic.Order(1))),
						MaxDifficulty: semantic.NewFloat(section, cfg.Key("mission_type"), semantic.Precision(2), semantic.OptsF(semantic.Order(2))),
						Weight:        semantic.NewInt(section, cfg.Key("mission_type"), semantic.Order(3)),
					}
					mission_type.Map(section)
					faction.MissionType = mission_type

					for index, _ := range section.ParamMap[cfg.Key("npc")] {
						faction.Npcs = append(faction.Npcs,
							semantic.NewString(mbase_section, cfg.Key("weight"), semantic.OptsS(semantic.Index(index))))
					}
					base.BaseFactions = append(base.BaseFactions, faction)
					base.BaseFactionsMap[faction.Faction.Get()] = faction
				case "[mroom]":
					mroom := &Mroom{
						Nickname:         semantic.NewString(section, cfg.Key("nickname"), semantic.WithLowercaseS(), semantic.WithoutSpacesS()),
						CharacterDensity: semantic.NewInt(section, cfg.Key("character_density")),
						Bartrender:       semantic.NewString(section, cfg.Key("fixture"), semantic.OptsS(semantic.Order(0), semantic.Optional())),
					}
					mroom.Map(section)
					if mroom.Nickname.Get() == "bar" {
						base.Bar = mroom
					}
				case "[gf_npc]":
					npc := &NPC{
						Nickname:    semantic.NewString(section, cfg.Key("nickname"), semantic.WithLowercaseS(), semantic.WithoutSpacesS()),
						Room:        semantic.NewString(section, cfg.Key("room"), semantic.OptsS(semantic.Optional())),
						Affiliation: semantic.NewString(section, cfg.Key("affiliation"), semantic.WithLowercaseS(), semantic.WithoutSpacesS()),
					}
					npc.Map(section)

					for index, _ := range section.ParamMap[cfg.Key("bribe")] {
						bribe := &Bribe{
							Faction: semantic.NewString(section, cfg.Key("bribe"), semantic.OptsS(semantic.Index(index)), semantic.WithLowercaseS(), semantic.WithoutSpacesS()),
						}
						bribe.Map(section)
						npc.Bribes = append(npc.Bribes, bribe)
					}
					for range section.ParamMap[cfg.Key("rumor")] {
						rumor := &Rumor{}
						rumor.Map(section)
						npc.Rumors = append(npc.Rumors, rumor)
					}
					for range section.ParamMap[cfg.Key("misc")] {
						misn := &Mission{}
						misn.Map(section)
						npc.Missions = append(npc.Missions, misn)
					}
					for range section.ParamMap[cfg.Key("know")] {
						know := &Know{}
						know.Map(section)
						npc.Knows = append(npc.Knows, know)
					}

					base.NPCs = append(base.NPCs, npc)
				}

			}
		}

	}

	return frelconfig
}

func (frelconfig *Config) Write() *file.File {
	// TODO BEWARE A BUG to fix.
	// if having here frelconfig.Render()
	// everything is still correct as typing
	// but the file is not getting written in darklint
	// This bug may be is going through my other code
	inifile := frelconfig.File.Render()
	inifile.Write(inifile.File)
	return inifile.File
}

type BaseChance struct {
	Base   string
	Chance float64
}

func FactionBribes(config *Config) map[string]map[string]float64 {
	// for faction, chance at certain base
	var faction_rephacks map[string]map[string]float64 = make(map[string]map[string]float64)

	for _, base := range config.Bases {

		// per faction chance at base
		logus.Log.Debug("base=" + base.Nickname.Get())
		var base_bribe_chances map[string]float64 = make(map[string]float64)
		var faction_members map[string]int = make(map[string]int)
		for _, npc := range base.NPCs {
			faction_members[npc.Affiliation.Get()] += 1

		}
		for _, npc := range base.NPCs {
			if base.Bar == nil {
				continue
			}
			npc_nickname := npc.Nickname.Get()
			bartrender := base.Bar.Bartrender.Get()
			if npc_nickname == bartrender {
				for _, bribe := range npc.Bribes {
					chance_increase := 1 / float64(len(npc.Bribes)+len(npc.Rumors)+len(npc.Missions)+len(npc.Knows))
					base_bribe_chances[bribe.Faction.Get()] += chance_increase
				}
			} else {
				for _, bribe := range npc.Bribes {
					var weight float64 = 0
					if faction, ok := base.BaseFactionsMap[npc.Affiliation.Get()]; ok {
						weight = float64(faction.Weight.Get())

						if value, ok := faction_members[npc.Affiliation.Get()]; ok {
							if value != 0 {
								weight = weight / float64(value)
							}
						}
					}

					chance_increase := float64(weight/100) * 1 / float64(len(npc.Bribes)+len(npc.Rumors)+len(npc.Missions)+len(npc.Knows))
					base_bribe_chances[bribe.Faction.Get()] += chance_increase
				}
			}
		}

		for faction, chance := range base_bribe_chances {
			_, ok := faction_rephacks[faction]
			if !ok {
				faction_rephacks[faction] = make(map[string]float64)
			}
			faction_rephacks[faction][base.Nickname.Get()] += chance

			if faction_rephacks[faction][base.Nickname.Get()] > 1.0 {
				faction_rephacks[faction][base.Nickname.Get()] = 1.0
			}
		}

	}
	return faction_rephacks
}
