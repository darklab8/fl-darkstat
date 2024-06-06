package configs_export

import (
	"fmt"
	"strings"

	"github.com/darklab8/fl-configs/configs/configs_mapped/freelancer_mapped/data_mapped/missions_mapped/mbases_mapped"
)

type Reputation struct {
	Name     string
	Rep      float64
	Empathy  float64
	Nickname string
}

type Faction struct {
	Name      string
	ShortName string
	Nickname  string

	ObjectDestruction float64
	MissionSuccess    float64
	MissionFailure    float64
	MissionAbort      float64

	InfonameID  int
	InfocardID  int
	Infocard    InfocardKey
	Reputations []Reputation
	Rephacks    []Rephack
}

type Rephack struct {
	BaseName   string
	BaseOwner  string
	BaseSystem string

	BaseNickname string
	Chance       float64
}

func (e *Exporter) GetFactions(bases []Base) []Faction {
	var factions []Faction = make([]Faction, 0, 100)

	var basemap map[string]Base = make(map[string]Base)
	for _, base := range bases {
		basemap[base.Nickname] = base
	}

	// for faction, at base, chance
	faction_rephacks := mbases_mapped.FactionRephacks(e.configs.MBases)

	for _, group := range e.configs.InitialWorld.Groups {
		var nickname string = group.Nickname.Get()
		faction := Faction{
			Nickname:   nickname,
			InfonameID: group.IdsName.Get(),
			InfocardID: group.IdsInfo.Get(),
			Infocard:   InfocardKey(nickname),
		}

		if rephacks, ok := faction_rephacks[nickname]; ok {

			for base, chance := range rephacks {
				rephack := Rephack{
					BaseNickname: base,
					Chance:       chance,
				}

				if base_info, ok := basemap[base]; ok {
					rephack.BaseName = base_info.Name
					rephack.BaseOwner = base_info.FactionName
					rephack.BaseSystem = base_info.System
				}

				faction.Rephacks = append(faction.Rephacks, rephack)
			}
		}
		faction.Name = e.GetInfocardName(group.IdsName.Get(), faction.Nickname)

		e.exportInfocards(InfocardKey(nickname), group.IdsInfo.Get())

		faction.ShortName = e.GetInfocardName(group.IdsShortName.Get(), faction.Nickname)

		empathy_rates, empathy_exists := e.configs.Empathy.RepoChangeMap[faction.Nickname]

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

			target_faction := e.configs.InitialWorld.GroupsMap[rep_to_add.Nickname]
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
