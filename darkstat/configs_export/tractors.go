package configs_export

import (
	"fmt"
	"sort"
	"strings"

	"github.com/darklab8/fl-darkstat/configs/cfg"
	"github.com/darklab8/fl-darkstat/configs/configs_mapped/freelancer_mapped/data_mapped/initialworld/flhash"
	"github.com/darklab8/fl-darkstat/configs/configs_settings/logus"
	"github.com/darklab8/fl-darkstat/configs/discovery/playercntl_rephacks"
	"github.com/darklab8/go-typelog/typelog"
)

type Rephack struct {
	FactionName string                      `json:"faction_name"`
	FactionNick cfg.FactionNick             `json:"faction_nickname"`
	Reputation  float64                     `json:"reputation"`
	RepType     playercntl_rephacks.RepType `json:"rep_type"`
}

type DiscoveryIDRephacks struct {
	Rephacks map[cfg.FactionNick]Rephack `json:"rephacks"`
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
	Name       string `json:"name"`
	Price      int    `json:"price"`
	MaxLength  int    `json:"max_length"`
	ReachSpeed int    `json:"reach_speed"`

	Lootable      bool            `json:"lootable"`
	Nickname      cfg.TractorID   `json:"nickname"`
	NicknameHash  flhash.HashCode `json:"nickname_hash" format:"int64"`
	ShortNickname string          `json:"short_nickname"`
	NameID        int             `json:"name_id"`
	InfoID        int             `json:"info_id"`

	Bases map[cfg.BaseUniNick]*MarketGood `json:"-" swaggerignore:"true"`
	DiscoveryIDRephacks
	Mass float64 `json:"mass"`
}

func (e *Exporter) GetFactionName(nickname cfg.FactionNick) string {
	if group, ok := e.Configs.InitialWorld.GroupsMap[string(nickname)]; ok {
		return e.GetInfocardName(group.IdsName.Get(), string(nickname))
	}
	return ""
}

func (e *Exporter) GetTractors() []*Tractor {
	var tractors []*Tractor

	for tractor_id, tractor_info := range e.Configs.Equip.Tractors {
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

		tractor.NicknameHash = flhash.HashNickname(string(tractor.Nickname))
		e.Hashes[string(tractor.Nickname)] = tractor.NicknameHash

		tractor.MaxLength = tractor_info.MaxLength.Get()
		tractor.ReachSpeed = tractor_info.ReachSpeed.Get()
		tractor.Lootable = tractor_info.Lootable.Get()
		tractor.NameID, _ = tractor_info.IdsName.GetValue()
		tractor.InfoID, _ = tractor_info.IdsInfo.GetValue()

		if good_info, ok := e.Configs.Goods.GoodsMap[string(tractor.Nickname)]; ok {
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

		if e.Configs.Discovery != nil {

			for faction_nick, faction := range e.Configs.Discovery.PlayercntlRephacks.DefaultReps {
				tractor.Rephacks[faction_nick] = Rephack{
					Reputation:  faction.Rep.Get(),
					RepType:     faction.GetRepType(),
					FactionNick: faction_nick,
					FactionName: e.GetFactionName(faction_nick),
				}
			}

			if faction, ok := e.Configs.Discovery.PlayercntlRephacks.RephacksByID[tractor.Nickname]; ok {

				if inherited_id, ok := faction.Inherits.GetValue(); ok {
					if faction, ok := e.Configs.Discovery.PlayercntlRephacks.RephacksByID[cfg.TractorID(inherited_id)]; ok {
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
