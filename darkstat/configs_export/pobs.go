package configs_export

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/darklab8/fl-darkstat/configs/cfg"
	"github.com/darklab8/fl-darkstat/configs/configs_mapped/freelancer_mapped/data_mapped/equipment_mapped"
	"github.com/darklab8/fl-darkstat/configs/configs_mapped/freelancer_mapped/data_mapped/equipment_mapped/equip_mapped"
	"github.com/darklab8/fl-darkstat/configs/configs_mapped/freelancer_mapped/data_mapped/initialworld"
	"github.com/darklab8/fl-darkstat/configs/configs_mapped/freelancer_mapped/data_mapped/initialworld/flhash"
	"github.com/darklab8/fl-darkstat/configs/configs_mapped/freelancer_mapped/data_mapped/universe_mapped"
	"github.com/darklab8/fl-darkstat/configs/configs_mapped/freelancer_mapped/infocard_mapped/infocard"
	"github.com/darklab8/fl-darkstat/configs/configs_settings/logus"
	"github.com/darklab8/fl-darkstat/configs/discovery/pob_goods"
	"github.com/darklab8/go-typelog/typelog"
	"github.com/darklab8/go-utils/utils/ptr"
)

type ShopItem struct {
	pob_goods.ShopItem
	Nickname string `json:"nickname"`
	Name     string `json:"name"`
	Category string `json:"category"`
}

type DefenseMode int

func (d DefenseMode) ToStr() string {
	switch d {

	case 1:
		return "SRP Whitelist > Blacklist > IFF Standing, Anyone with good standing"
	case 2:
		return "Whitelist > Nodock, Whitelisted ships only"
	case 3:
		return "Whitelist > Hostile, Whitelisted ships only"
	default:
		return "not recognized"
	}
}

type PoBWithoutGoods struct {
	Nickname string `json:"nickname"`
	Name     string `json:"name"`

	Pos         *string      `json:"pos"`
	Level       *int         `json:"level"`
	Money       *int         `json:"money"`
	Health      *float64     `json:"health"`
	DefenseMode *DefenseMode `json:"defense_mode"`

	SystemNick  *string `json:"system_nickname"`
	SystemName  *string `json:"system_name"` // SystemHash      *flhash.HashCode `json:"system"`      //: 2745655887,
	FactionNick *string `json:"faction_nickname"`
	FactionName *string `json:"faction_name"` // AffiliationHash *flhash.HashCode `json:"affiliation"` //: 2620,

	ForumThreadUrl *string `json:"forum_thread_url"`
	CargoSpaceLeft *int    `json:"cargospace"`

	BasePos     *cfg.Vector `json:"base_pos"`
	SectorCoord *string     `json:"sector_coord"`
	Region      *string     `json:"region_name"`
}

// also known as Player Base Station
type PoB struct {
	PoBWithoutGoods
	ShopItems []*ShopItem `json:"shop_items"`
}

func (b PoB) GetNickname() string { return string(b.Nickname) }

type PoBGood struct {
	Nickname              string `json:"nickname"`
	Name                  string `json:"name"`
	TotalBuyableFromBases int    `json:"total_buyable_from_bases"`
	TotalSellableToBases  int    `json:"total_sellable_to_bases"`

	BestPriceToBuy  *int `json:"best_price_to_buy"`
	BestPriceToSell *int `json:"best_price_to_sell"`

	Category string         `json:"category"`
	Bases    []*PoBGoodBase `json:"bases"`

	AnyBaseSells bool `json:"any_base_sells"`
	AnyBaseBuys  bool `json:"any_base_buys"`
}

func (b PoBGood) GetNickname() string { return string(b.Nickname) }

func (good PoBGood) BaseSells() bool { return good.AnyBaseSells }
func (good PoBGood) BaseBuys() bool  { return good.AnyBaseBuys }

type PoBGoodBase struct {
	ShopItem *ShopItem        `json:"shop_item"`
	Base     *PoBWithoutGoods `json:"base"`
}

// Exporting only with position ones
func (e *Exporter) PoBsToBases(pobs []*PoB) []*Base {
	var bases []*Base

	for _, pob := range pobs {
		if pob.BasePos == nil {
			continue
		}

		base := &Base{
			Nickname:           cfg.BaseUniNick(pob.Nickname),
			NicknameHash:       flhash.HashNickname(pob.Nickname),
			Name:               fmt.Sprintf("(PoB) %s", pob.Name),
			Pos:                *pob.BasePos,
			System:             *pob.SystemName,
			SystemNickname:     *pob.SystemNick,
			Region:             *pob.Region,
			SectorCoord:        *pob.SectorCoord,
			MarketGoodsPerNick: map[CommodityKey]*MarketGood{},
			IsPob:              true,
		}
		if pob.FactionName != nil {
			base.FactionName = *pob.FactionName
		}
		bases = append(bases, base)

		for _, pob_good := range pob.ShopItems {
			market_good := &MarketGood{
				GoodInfo:             e.GetGoodInfo(pob_good.Nickname),
				IsServerSideOverride: true,
				ShipClass:            -1,
			}
			if pob_good.BaseBuys() {
				market_good.PriceBaseBuysFor = ptr.Ptr(pob_good.SellPrice)
			}
			if pob_good.BaseSells() {
				market_good.PriceBaseSellsFor = pob_good.Price
				market_good.BaseSells = true
			}
			if market_good.Category == "commodity" {
				equipment := e.Configs.Equip.CommoditiesMap[market_good.Nickname]
				for _, volume := range equipment.Volumes {
					volumed_good := market_good
					volumed_good.Volume = volume.Volume.Get()
					volumed_good.ShipClass = volume.GetShipClass()
					base.MarketGoodsPerNick[GetCommodityKey(volumed_good.Nickname, volumed_good.ShipClass)] = volumed_good
				}
			} else {
				base.MarketGoodsPerNick[GetCommodityKey(market_good.Nickname, market_good.ShipClass)] = market_good
			}
		}
	}
	return bases
}

func (e *Exporter) GetPoBGoods(pobs []*PoB) []*PoBGood {
	pobs_goods_by_nick := make(map[string]*PoBGood)
	var pob_goods []*PoBGood

	for _, pob := range pobs {
		for _, good := range pob.ShopItems {
			pob_good, found_good := pobs_goods_by_nick[good.Nickname]
			if !found_good {
				pob_good = &PoBGood{
					Nickname: good.Nickname,
					Name:     good.Name,
					Category: good.Category,
				}
				pobs_goods_by_nick[good.Nickname] = pob_good
			}
			pob_good.Bases = append(pob_good.Bases, &PoBGoodBase{Base: &pob.PoBWithoutGoods, ShopItem: good})
		}
	}

	for _, item := range pobs_goods_by_nick {
		for _, pob := range item.Bases {
			if pob.ShopItem.BaseSells() {
				item.AnyBaseSells = true
				item.TotalBuyableFromBases += pob.ShopItem.Quantity - pob.ShopItem.MinStock

				if item.BestPriceToBuy == nil {
					item.BestPriceToBuy = ptr.Ptr(pob.ShopItem.Price)
				}
				if pob.ShopItem.Price < *item.BestPriceToBuy {
					item.BestPriceToBuy = ptr.Ptr(pob.ShopItem.Price)
				}
			}
			if pob.ShopItem.BaseBuys() {
				item.AnyBaseBuys = true
				sellable_to_current_base := pob.ShopItem.MaxStock - pob.ShopItem.Quantity

				if pob.Base.CargoSpaceLeft != nil {
					if *pob.Base.CargoSpaceLeft < sellable_to_current_base {
						sellable_to_current_base = *pob.Base.CargoSpaceLeft
					}
				}

				item.TotalSellableToBases += sellable_to_current_base

				if item.BestPriceToSell == nil {
					item.BestPriceToSell = ptr.Ptr(pob.ShopItem.SellPrice)
				}
				if pob.ShopItem.SellPrice > *item.BestPriceToSell {
					item.BestPriceToSell = ptr.Ptr(pob.ShopItem.SellPrice)
				}
			}
		}

		pob_goods = append(pob_goods, item)
	}

	return pob_goods
}

func (e *Exporter) GetPoBs() []*PoB {
	var pobs []*PoB

	systems_by_hash := make(map[flhash.HashCode]*universe_mapped.System)
	factions_by_hash := make(map[flhash.HashCode]*initialworld.Group)
	for _, system_info := range e.Configs.Universe.Systems {
		nickname := system_info.Nickname.Get()
		system_hash := flhash.HashNickname(nickname)
		systems_by_hash[system_hash] = system_info
	}
	for _, group_info := range e.Configs.InitialWorld.Groups {
		nickname := group_info.Nickname.Get()
		group_hash := flhash.HashFaction(nickname)
		factions_by_hash[group_hash] = group_info
	}
	goods_by_hash := make(map[flhash.HashCode]*equip_mapped.Item)
	for _, item := range e.Configs.Equip.Items {
		nickname := item.Nickname.Get()
		hash := flhash.HashNickname(nickname)
		goods_by_hash[hash] = item
		e.exportInfocards(InfocardKey(nickname), item.IdsInfo.Get())
	}

	ships_by_hash := make(map[flhash.HashCode]*equipment_mapped.Ship)
	for _, item := range e.Configs.Goods.Ships {
		nickname := item.Nickname.Get()
		hash := flhash.HashNickname(nickname)
		ships_by_hash[hash] = item
	}

	for _, pob_info := range e.Configs.Discovery.PlayerOwnedBases.Bases {

		var pob *PoB = &PoB{
			PoBWithoutGoods: PoBWithoutGoods{
				Nickname:       pob_info.Nickname,
				Name:           pob_info.Name,
				Pos:            pob_info.Pos,
				Level:          pob_info.Level,
				Money:          pob_info.Money,
				Health:         pob_info.Health,
				CargoSpaceLeft: pob_info.CargoSpaceLeft,
			},
		}
		if pob_info.DefenseMode != nil {
			pob.DefenseMode = (*DefenseMode)(pob_info.DefenseMode)
		}
		if pob_info.Pos != nil {
			pob.BasePos = StrPosToVectorPos(*pob_info.Pos)
		}

		pob.ForumThreadUrl = pob_info.ForumThreadUrl
		if pob_info.SystemHash != nil {
			if system, ok := systems_by_hash[*pob_info.SystemHash]; ok {
				pob.SystemNick = ptr.Ptr(system.Nickname.Get())
				pob.SystemName = ptr.Ptr(e.GetInfocardName(system.StridName.Get(), system.Nickname.Get()))

				pob.Region = ptr.Ptr(e.GetRegionName(system))
				if pob.BasePos != nil {
					pob.SectorCoord = ptr.Ptr(VectorToSectorCoord(system, *pob.BasePos))
				}
			}
		}

		if pob_info.AffiliationHash != nil {
			if faction, ok := factions_by_hash[*pob_info.AffiliationHash]; ok {
				pob.FactionNick = ptr.Ptr(faction.Nickname.Get())
				pob.FactionName = ptr.Ptr(e.GetInfocardName(faction.IdsName.Get(), faction.Nickname.Get()))
			}
		}

		for _, shop_item := range pob_info.ShopItems {
			good := &ShopItem{ShopItem: shop_item}

			if item, ok := goods_by_hash[flhash.HashCode(shop_item.Id)]; ok {
				good.Nickname = item.Nickname.Get()
				good.Name = e.GetInfocardName(item.IdsName.Get(), item.Nickname.Get())
				good.Category = item.Category
			} else {
				if ship, ok := ships_by_hash[flhash.HashCode(shop_item.Id)]; ok {
					ship_hull := e.Configs.Goods.ShipHullsMap[ship.Hull.Get()]
					ship_nickname := ship_hull.Ship.Get()
					shiparch := e.Configs.Shiparch.ShipsMap[ship_nickname]
					good.Nickname = ship_nickname
					good.Category = "ship"
					good.Name = e.GetInfocardName(shiparch.IdsName.Get(), ship_nickname)
				} else {
					logus.Log.Warn("unidentified shop item", typelog.Any("shop_item.Id", shop_item.Id))
				}
			}

			pob.ShopItems = append(pob.ShopItems, good)
		}

		var sb InfocardBuilder
		sb.WriteLineStr(pob.Name)
		sb.WriteLineStr("")

		if pob_info.Pos == nil && len(pob_info.InfocardParagraphs) == 0 {
			sb.WriteLine(InfocardPhrase{Phrase: "infocard:", Bold: true})
			sb.WriteLineStr("no access (toggle pos permission in pob account manager)")
			sb.WriteLineStr("")
		}

		for _, paragraph := range pob_info.InfocardParagraphs {
			sb.WriteLineStr(paragraph)
			sb.WriteLineStr("")
		}

		if pob_info.DefenseMode != nil {
			sb.WriteLine(InfocardPhrase{Phrase: "Defense mode:", Bold: true})
			sb.WriteLineStr((*DefenseMode)(pob_info.DefenseMode).ToStr())
			sb.WriteLineStr("")
		} else {
			sb.WriteLine(InfocardPhrase{Phrase: "docking permissions:", Bold: true})
			sb.WriteLineStr("no access (toggle defense mode in pob account manager)")
			sb.WriteLineStr("")
		}
		if len(pob_info.SrpFactionHashList) > 0 || len(pob_info.SrpTagList) > 0 || len(pob_info.SrpNameList) > 0 {
			sb.WriteLine(InfocardPhrase{Phrase: "Docking allias(srp,ignore rep):", Bold: true})
			sb.WriteLineStr(e.fmt_factions_to_str(factions_by_hash, pob_info.SrpFactionHashList))
			sb.WriteLineStr(fmt.Sprintf("tags: %s", fmt_docking_tags(pob_info.SrpTagList)))
			sb.WriteLineStr(fmt.Sprintf("names: %s", fmt_docking_tags(pob_info.SrpNameList)))
			sb.WriteLineStr("")
		}

		if len(pob_info.AllyFactionHashList) > 0 || len(pob_info.AllyTagList) > 0 || len(pob_info.AllyNameList) > 0 {
			sb.WriteLine(InfocardPhrase{Phrase: "Docking allias(IFF rep still affects):", Bold: true})
			sb.WriteLineStr(e.fmt_factions_to_str(factions_by_hash, pob_info.AllyFactionHashList))
			sb.WriteLineStr(fmt.Sprintf("tags: %s", fmt_docking_tags(pob_info.AllyTagList)))
			sb.WriteLineStr(fmt.Sprintf("names: %s", fmt_docking_tags(pob_info.AllyNameList)))
			sb.WriteLineStr("")
		}

		if len(pob_info.HostileFactionHashList) > 0 || len(pob_info.HostileTagList) > 0 || len(pob_info.HostileNameList) > 0 {
			sb.WriteLine(InfocardPhrase{Phrase: "Docking enemies:", Bold: true})
			sb.WriteLineStr(e.fmt_factions_to_str(factions_by_hash, pob_info.HostileFactionHashList))
			sb.WriteLineStr(fmt.Sprintf("tags: %s", fmt_docking_tags(pob_info.HostileTagList)))
			sb.WriteLineStr(fmt.Sprintf("names: %s", fmt_docking_tags(pob_info.HostileNameList)))
			sb.WriteLineStr("")
		}

		// TODO add pob infocards here
		e.Infocards[InfocardKey(pob.Nickname)] = sb.Lines
		e.Configs.Infocards.Infonames[int(flhash.HashNickname(pob.Nickname))] = infocard.Infoname(pob.Name)

		pobs = append(pobs, pob)
	}
	return pobs
}

type PobShopItem struct {
	*ShopItem
	PoBName     string
	PobNickname string

	System      *universe_mapped.System
	SystemNick  string
	SystemName  string
	FactionNick string
	FactionName string
	BasePos     *cfg.Vector
}

func (e *Exporter) get_pob_buyable() map[string][]*PobShopItem {
	if e.pob_buyable_cache != nil {
		return e.pob_buyable_cache
	}

	e.pob_buyable_cache = make(map[string][]*PobShopItem)

	// TODO refactor copy repeated code may be
	systems_by_hash := make(map[flhash.HashCode]*universe_mapped.System)
	factions_by_hash := make(map[flhash.HashCode]*initialworld.Group)
	for _, system_info := range e.Configs.Universe.Systems {
		nickname := system_info.Nickname.Get()
		system_hash := flhash.HashNickname(nickname)
		systems_by_hash[system_hash] = system_info
	}
	for _, group_info := range e.Configs.InitialWorld.Groups {
		nickname := group_info.Nickname.Get()
		group_hash := flhash.HashFaction(nickname)
		factions_by_hash[group_hash] = group_info
	}
	goods_by_hash := make(map[flhash.HashCode]*equip_mapped.Item)
	for _, item := range e.Configs.Equip.Items {
		nickname := item.Nickname.Get()
		hash := flhash.HashNickname(nickname)
		goods_by_hash[hash] = item
		e.exportInfocards(InfocardKey(nickname), item.IdsInfo.Get())
	}
	ships_by_hash := make(map[flhash.HashCode]*equipment_mapped.Ship)
	for _, item := range e.Configs.Goods.Ships {
		nickname := item.Nickname.Get()
		hash := flhash.HashNickname(nickname)
		ships_by_hash[hash] = item
	}

	for _, pob_info := range e.Configs.Discovery.PlayerOwnedBases.Bases {
		for _, shop_item := range pob_info.ShopItems {
			var good *ShopItem = &ShopItem{ShopItem: shop_item}
			if item, ok := goods_by_hash[flhash.HashCode(shop_item.Id)]; ok {
				good.Nickname = item.Nickname.Get()
				good.Name = e.GetInfocardName(item.IdsName.Get(), item.Nickname.Get())
				good.Category = item.Category
			} else {
				if ship, ok := ships_by_hash[flhash.HashCode(shop_item.Id)]; ok {
					ship_hull := e.Configs.Goods.ShipHullsMap[ship.Hull.Get()]
					ship_nickname := ship_hull.Ship.Get()
					shiparch := e.Configs.Shiparch.ShipsMap[ship_nickname]
					good.Nickname = ship_nickname
					good.Category = "ship"
					good.Name = e.GetInfocardName(shiparch.IdsName.Get(), ship_nickname)
				} else {
					logus.Log.Warn("unidentified shop item", typelog.Any("shop_item.Id", shop_item.Id))
				}
			}
			pob_item := &PobShopItem{
				ShopItem:    good,
				PobNickname: pob_info.Nickname,
				PoBName:     pob_info.Name,
			}

			if pob_info.SystemHash != nil {
				if system, ok := systems_by_hash[*pob_info.SystemHash]; ok {
					pob_item.SystemNick = system.Nickname.Get()
					pob_item.SystemName = e.GetInfocardName(system.StridName.Get(), system.Nickname.Get())
					pob_item.System = system
				}
			}
			if pob_info.AffiliationHash != nil {
				if faction, ok := factions_by_hash[*pob_info.AffiliationHash]; ok {
					pob_item.FactionNick = faction.Nickname.Get()
					pob_item.FactionName = e.GetInfocardName(faction.IdsName.Get(), faction.Nickname.Get())
				}
			}

			if pob_info.Pos != nil {
				pob_item.BasePos = StrPosToVectorPos(*pob_info.Pos)
			}

			e.pob_buyable_cache[good.Nickname] = append(e.pob_buyable_cache[good.Nickname], pob_item)
		}
	}
	return e.pob_buyable_cache
}

func (e *Exporter) fmt_factions_to_str(factions_by_hash map[flhash.HashCode]*initialworld.Group, faction_hashes []*flhash.HashCode) string {
	var sb strings.Builder

	sb.WriteString("factions: [")

	for index, faction_hash := range faction_hashes {
		if faction, ok := factions_by_hash[*faction_hash]; ok {
			sb.WriteString(e.GetInfocardName(faction.IdsName.Get(), faction.Nickname.Get()))
			if index != len(faction_hashes)-1 {
				sb.WriteString(", ")
			}
		} else {
			logus.Log.Error("faction hash is invalid", typelog.Any("hash", *faction_hash))
		}
	}
	sb.WriteString("]")
	return sb.String()
}

func fmt_docking_tags(tags_or_names []string) string {
	return fmt.Sprintf("[%s]", strings.Join(tags_or_names, ", "))
}

func StrPosToVectorPos(value string) *cfg.Vector {
	coords := strings.Split(value, ",")
	x, err1 := strconv.ParseFloat(strings.ReplaceAll(coords[0], " ", ""), 64)
	y, err2 := strconv.ParseFloat(strings.ReplaceAll(coords[1], " ", ""), 64)
	z, err3 := strconv.ParseFloat(strings.ReplaceAll(coords[2], " ", ""), 64)
	logus.Log.CheckPanic(err1, "failed parsing x coord", typelog.Any("pos", value))
	logus.Log.CheckPanic(err2, "failed parsing y coord", typelog.Any("pos", value))
	logus.Log.CheckPanic(err3, "failed parsing z coord", typelog.Any("pos", value))

	return &cfg.Vector{X: x, Y: y, Z: z}
}
