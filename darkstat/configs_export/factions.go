package configs_export

import (
	"fmt"
	"strings"

	"github.com/darklab8/fl-darkstat/configs/cfg"
	"github.com/darklab8/fl-darkstat/configs/configs_mapped/freelancer_mapped/data_mapped/missions_mapped/mbases_mapped"
	"github.com/darklab8/fl-darkstat/configs/configs_mapped/freelancer_mapped/data_mapped/universe_mapped"
	"github.com/darklab8/fl-darkstat/darkstat/configs_export/infocarder"
)

type Reputation struct {
	Name     string  `json:"name"  validate:"required"`
	Rep      float64 `json:"rep"  validate:"required"`
	Empathy  float64 `json:"empathy"  validate:"required"`
	Nickname string  `json:"nickname"  validate:"required"`
}

type Faction struct {
	Name      string `json:"name"  validate:"required"`
	ShortName string `json:"short_name"  validate:"required"`
	Nickname  string `json:"nickname"  validate:"required"`

	ObjectDestruction float64 `json:"object_destruction"  validate:"required"`
	MissionSuccess    float64 `json:"mission_success"  validate:"required"`
	MissionFailure    float64 `json:"mission_failure"  validate:"required"`
	MissionAbort      float64 `json:"mission_abort"  validate:"required"`

	InfonameID  int                    `json:"infoname_id"  validate:"required"`
	InfocardID  int                    `json:"infocard_id"  validate:"required"`
	InfocardKey infocarder.InfocardKey `json:"-" swaggerignore:"true"`
	Reputations []Reputation           `json:"reputations"  validate:"required"`
	Bribes      []Bribe                `json:"bribe"  validate:"required"`
}

func (b Faction) GetNickname() string { return string(b.Nickname) }

type Bribe struct {
	BaseNickname string `json:"base_nickname"  validate:"required"`
	BaseInfo
	Chance float64 `json:"chance"  validate:"required"`
}

func (e *Exporter) GetFactions(bases []*Base) []Faction {
	var factions []Faction = make([]Faction, 0, 100)

	var basemap map[cfg.BaseUniNick]*Base = make(map[cfg.BaseUniNick]*Base)
	for _, base := range bases {
		basemap[base.Nickname] = base
	}

	// for faction, at base, chance
	faction_rephacks := mbases_mapped.FactionBribes(e.Mapped.MBases)

	for _, group := range e.Mapped.InitialWorld.Groups {
		var nickname string = group.Nickname.Get()
		faction := Faction{
			Nickname:    nickname,
			InfonameID:  group.IdsName.Get(),
			InfocardID:  group.IdsInfo.Get(),
			InfocardKey: infocarder.InfocardKey(nickname),
		}

		if rephacks, ok := faction_rephacks[nickname]; ok {

			for base, chance := range rephacks {
				rephack := Bribe{
					BaseNickname: base,
					Chance:       chance,
					BaseInfo:     e.GetBaseInfo(universe_mapped.BaseNickname(base)),
				}

				faction.Bribes = append(faction.Bribes, rephack)
			}
		}
		faction.Name = e.GetInfocardName(group.IdsName.Get(), faction.Nickname)

		e.exportInfocards(infocarder.InfocardKey(nickname), group.IdsInfo.Get())

		faction.ShortName = e.GetInfocardName(group.IdsShortName.Get(), faction.Nickname)

		empathy_rates, empathy_exists := e.Mapped.Empathy.RepoChangeMap[faction.Nickname]

		if empathy_exists {
			faction.ObjectDestruction = empathy_rates.ObjectDestruction.Get()
			faction.MissionSuccess = empathy_rates.MissionSuccess.Get()
			faction.MissionFailure = empathy_rates.MissionFailure.Get()
			faction.MissionAbort = empathy_rates.MissionAbort.Get()
		}

		for _, reputation := range group.Relationships {
			rep_to_add := &Reputation{}
			rep_to_add.Nickname = reputation.TargetNickname.Get()
			rep_to_add.Rep = reputation.Rep.Get()

			target_faction := e.Mapped.InitialWorld.GroupsMap[rep_to_add.Nickname]
			defer func() {
				if r := recover(); r != nil {
					fmt.Println("Recovered in f", r)
					fmt.Println("recovered rep_to_add.Nickname", rep_to_add.Nickname)
					panic(r)
				}
			}()

			if target_faction != nil {
				rep_to_add.Name = e.GetInfocardName(target_faction.IdsName.Get(), rep_to_add.Nickname)
			}

			if empathy_exists {
				if empathy_rate, ok := empathy_rates.EmpathyRatesMap[rep_to_add.Nickname]; ok {
					rep_to_add.Empathy = empathy_rate.RepoChange.Get()
				}
			}

			faction.Reputations = append(faction.Reputations, *rep_to_add)
		}

		factions = append(factions, faction)

	}

	return factions
}

func FilterToUsefulFactions(factions []Faction) []Faction {
	var useful_factions []Faction = make([]Faction, 0, len(factions))
	for _, item := range factions {
		if Empty(item.Name) || strings.Contains(item.Name, "_grp") {
			continue
		}

		useful_factions = append(useful_factions, item)
	}
	return useful_factions
}

func FilterToUsefulBribes(factions []Faction) []Faction {
	var useful_factions []Faction = make([]Faction, 0, len(factions))
	for _, item := range factions {
		if len(item.Bribes) == 0 {
			continue
		}

		useful_factions = append(useful_factions, item)
	}
	return useful_factions
}
