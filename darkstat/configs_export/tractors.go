package configs_export

import (
	"fmt"
	"sort"
	"strings"

	"github.com/darklab8/fl-configs/configs/cfgtype"
	"github.com/darklab8/fl-configs/configs/configs_mapped/freelancer_mapped/data_mapped/initialworld/flhash"
	"github.com/darklab8/fl-configs/configs/discovery/playercntl_rephacks"
)

type Rephack struct {
	FactionName string
	FactionNick cfgtype.FactionNick
	Reputation  float64
	RepType     playercntl_rephacks.RepType
}

type DiscoveryIDRephacks struct {
	Rephacks map[cfgtype.FactionNick]Rephack
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
	Name       string
	Price      int
	MaxLength  int
	ReachSpeed int

	Lootable      bool
	Nickname      cfgtype.TractorID
	NicknameHash  flhash.HashCode
	ShortNickname string
	NameID        int
	InfoID        int

	Bases map[cfgtype.BaseUniNick]*GoodAtBase
	DiscoveryIDRephacks
	Mass float64
}

func (e *Exporter) GetFactionName(nickname cfgtype.FactionNick) string {
	if group, ok := e.Configs.InitialWorld.GroupsMap[string(nickname)]; ok {
		return e.GetInfocardName(group.IdsName.Get(), string(nickname))
	}
	return ""
}

func (e *Exporter) GetTractors() []Tractor {
	var tractors []Tractor

	for tractor_id, tractor_info := range e.Configs.Equip.Tractors {
		tractor := Tractor{
			ShortNickname: fmt.Sprintf("i%d", tractor_id),
			DiscoveryIDRephacks: DiscoveryIDRephacks{
				Rephacks: make(map[cfgtype.FactionNick]Rephack),
			},
			Bases: make(map[cfgtype.BaseUniNick]*GoodAtBase),
		}
		tractor.Mass, _ = tractor_info.Mass.GetValue()

		tractor.Nickname = cfgtype.TractorID(tractor_info.Nickname.Get())
		tractor.NicknameHash = flhash.HashNickname(string(tractor.Nickname))
		e.Hashes[string(tractor.Nickname)] = tractor.NicknameHash

		tractor.MaxLength = tractor_info.MaxLength.Get()
		tractor.ReachSpeed = tractor_info.ReachSpeed.Get()
		tractor.Lootable = tractor_info.Lootable.Get()
		tractor.NameID = tractor_info.IdsName.Get()
		tractor.InfoID = tractor_info.IdsInfo.Get()

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
					if faction, ok := e.Configs.Discovery.PlayercntlRephacks.RephacksByID[cfgtype.TractorID(inherited_id)]; ok {
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

func (e *Exporter) FilterToUsefulTractors(tractors []Tractor) []Tractor {
	var buyable_tractors []Tractor = make([]Tractor, 0, len(tractors))
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
