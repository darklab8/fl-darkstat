package configs_export

import (
	"fmt"
	"sort"
	"strings"

	"github.com/darklab8/fl-darkstat/configs/cfg"
	"github.com/darklab8/fl-darkstat/configs/configs_settings/logus"
	"github.com/darklab8/fl-darkstat/configs/discovery/playercntl_rephacks"
	"github.com/darklab8/go-typelog/typelog"
)

type Rephack struct {
	FactionName string                      `json:"faction_name" validate:"required"`
	FactionNick cfg.FactionNick             `json:"faction_nickname" validate:"required"`
	Reputation  float64                     `json:"reputation" validate:"required"`
	RepType     playercntl_rephacks.RepType `json:"rep_type" validate:"required"`
}

type DiscoveryIDRephacks struct {
	Rephacks map[cfg.FactionNick]Rephack `json:"rephacks" validate:"required"`
}

func (r DiscoveryIDRephacks) GetRephacksList() []Rephack {

	var result []Rephack
	for _, rephack := range r.Rephacks {

		result = append(result, rephack)
	}
	sort.Slice(result, func(i, j int) bool {
		return result[i].Reputation > result[j].Reputation
	})
	return result

}

type Tractor struct {
	Name       string `json:"name" validate:"required"`
	Price      int    `json:"price" validate:"required"`
	MaxLength  int    `json:"max_length" validate:"required"`
	ReachSpeed int    `json:"reach_speed" validate:"required"`

	Lootable      bool          `json:"lootable" validate:"required"`
	Nickname      cfg.TractorID `json:"nickname" validate:"required"`
	ShortNickname string        `json:"short_nickname" validate:"required"`
	NameID        int           `json:"name_id" validate:"required"`
	InfoID        int           `json:"info_id" validate:"required"`

	Bases map[cfg.BaseUniNick]*MarketGood `json:"-" swaggerignore:"true"`
	DiscoveryIDRephacks
	Mass float64 `json:"mass" validate:"required"`
}

func (e *Exporter) GetFactionName(nickname cfg.FactionNick) string {
	if group, ok := e.mapped.InitialWorld.GroupsMap[string(nickname)]; ok {
		return e.GetInfocardName(group.IdsName.Get(), string(nickname))
	}
	return ""
}

func (e *Exporter) GetTractors() []*Tractor {
	var tractors []*Tractor

	for tractor_id, tractor_info := range e.mapped.Equip.Tractors {
		tractor := &Tractor{
			Nickname:      cfg.TractorID(tractor_info.Nickname.Get()),
			ShortNickname: fmt.Sprintf("i%d", tractor_id),
			DiscoveryIDRephacks: DiscoveryIDRephacks{
				Rephacks: make(map[cfg.FactionNick]Rephack),
			},
			Bases: make(map[cfg.BaseUniNick]*MarketGood),
		}

		if _, ok := tractor_info.IdsName.GetValue(); !ok {
			logus.Log.Warn("tractor is not having defined ids_name", typelog.Any("nickname", tractor.Nickname))
		}
		tractor.Mass, _ = tractor_info.Mass.GetValue()
		tractor.MaxLength = tractor_info.MaxLength.Get()
		tractor.ReachSpeed = tractor_info.ReachSpeed.Get()
		tractor.Lootable = tractor_info.Lootable.Get()
		tractor.NameID, _ = tractor_info.IdsName.GetValue()
		tractor.InfoID, _ = tractor_info.IdsInfo.GetValue()

		if good_info, ok := e.mapped.Goods.GoodsMap[string(tractor.Nickname)]; ok {
			if price, ok := good_info.Price.GetValue(); ok {
				tractor.Price = price
				tractor.Bases = e.GetAtBasesSold(GetCommodityAtBasesInput{
					Nickname: good_info.Nickname.Get(),
					Price:    price,
				})
			}
		}

		tractor.Name = e.GetInfocardName(tractor.NameID, string(tractor.Nickname))

		e.exportInfocards(InfocardKey(tractor.Nickname), tractor.InfoID)

		if e.mapped.Discovery != nil {

			for faction_nick, faction := range e.mapped.Discovery.PlayercntlRephacks.DefaultReps {
				tractor.Rephacks[faction_nick] = Rephack{
					Reputation:  faction.Rep.Get(),
					RepType:     faction.GetRepType(),
					FactionNick: faction_nick,
					FactionName: e.GetFactionName(faction_nick),
				}
			}

			if faction, ok := e.mapped.Discovery.PlayercntlRephacks.RephacksByID[tractor.Nickname]; ok {

				if inherited_id, ok := faction.Inherits.GetValue(); ok {
					if faction, ok := e.mapped.Discovery.PlayercntlRephacks.RephacksByID[cfg.TractorID(inherited_id)]; ok {
						for faction_nick, rep := range faction.Reps {
							tractor.Rephacks[faction_nick] = Rephack{
								Reputation:  rep.Rep.Get(),
								RepType:     rep.GetRepType(),
								FactionNick: faction_nick,
								FactionName: e.GetFactionName(faction_nick),
							}
						}
					}
				}

				for faction_nick, rep := range faction.Reps {
					tractor.Rephacks[faction_nick] = Rephack{
						Reputation:  rep.Rep.Get(),
						RepType:     rep.GetRepType(),
						FactionNick: faction_nick,
						FactionName: e.GetFactionName(faction_nick),
					}
				}
			}
		}
		tractors = append(tractors, tractor)
	}
	return tractors
}

func (b Tractor) GetNickname() string { return string(b.Nickname) }

func (b Tractor) GetBases() map[cfg.BaseUniNick]*MarketGood { return b.Bases }

func (e *Exporter) FilterToUsefulTractors(tractors []*Tractor) []*Tractor {
	var buyable_tractors []*Tractor = make([]*Tractor, 0, len(tractors))
	for _, item := range tractors {

		if !e.Buyable(item.Bases) && (strings.Contains(strings.ToLower(item.Name), "discontinued") ||
			strings.Contains(strings.ToLower(item.Name), "not in use") ||
			strings.Contains(strings.ToLower(item.Name), strings.ToLower("Special Operative ID")) ||
			strings.Contains(strings.ToLower(item.Name), strings.ToLower("SRP ID")) ||
			strings.Contains(strings.ToLower(item.Name), strings.ToLower("Unused"))) {
			continue
		}
		buyable_tractors = append(buyable_tractors, item)
	}
	return buyable_tractors
}
