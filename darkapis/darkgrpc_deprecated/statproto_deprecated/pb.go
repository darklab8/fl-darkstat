package statproto_deprecated

type Pos struct {
	X float64 `protobuf:"fixed64,1,opt,name=x,proto3" json:"x,omitempty"`
	Y float64 `protobuf:"fixed64,2,opt,name=y,proto3" json:"y,omitempty"`
	Z float64 `protobuf:"fixed64,3,opt,name=z,proto3" json:"z,omitempty"`
}

type GetGunsInput struct {
	// "insert 'true' if wish to include market goods under 'market goods' key or not. Such data can add a lot of extra weight"
	IncludeMarketGoods bool `protobuf:"varint,1,opt,name=include_market_goods,json=includeMarketGoods,proto3" json:"include_market_goods,omitempty"`
	IncludeTechCompat  bool `protobuf:"varint,2,opt,name=include_tech_compat,json=includeTechCompat,proto3" json:"include_tech_compat,omitempty"`
	// Apply filtering same as darkstat does by default for its tab. Usually means showing only items that can be bought/crafted/or found
	FilterToUseful       bool `protobuf:"varint,3,opt,name=filter_to_useful,json=filterToUseful,proto3" json:"filter_to_useful,omitempty"`
	IncludeDamageBonuses bool `protobuf:"varint,4,opt,name=include_damage_bonuses,json=includeDamageBonuses,proto3" json:"include_damage_bonuses,omitempty"`
	// filters by item nicknames
	FilterNicknames []string `protobuf:"bytes,5,rep,name=filter_nicknames,json=filterNicknames,proto3" json:"filter_nicknames,omitempty"`
}

type GetTractorsInput struct {
	// By default not outputing market goods in case u wish to save network
	IncludeMarketGoods bool `protobuf:"varint,1,opt,name=include_market_goods,json=includeMarketGoods,proto3" json:"include_market_goods,omitempty"`
	// Apply filtering same as darkstat does by default for its tab. Usually means showing only items that can be bought/crafted/or found
	FilterToUseful bool `protobuf:"varint,2,opt,name=filter_to_useful,json=filterToUseful,proto3" json:"filter_to_useful,omitempty"`
	// filters by item nicknames
	FilterNicknames []string `protobuf:"bytes,3,rep,name=filter_nicknames,json=filterNicknames,proto3" json:"filter_nicknames,omitempty"`
	IncludeRephacks bool     `protobuf:"varint,4,opt,name=include_rephacks,json=includeRephacks,proto3" json:"include_rephacks,omitempty"`
}

type GetFactionsInput struct {
	IncludeReputations bool `protobuf:"varint,1,opt,name=include_reputations,json=includeReputations,proto3" json:"include_reputations,omitempty"`
	IncludeBribes      bool `protobuf:"varint,2,opt,name=include_bribes,json=includeBribes,proto3" json:"include_bribes,omitempty"`
	FilterToUseful     bool `protobuf:"varint,3,opt,name=filter_to_useful,json=filterToUseful,proto3" json:"filter_to_useful,omitempty"`
}

type GetEquipmentInput struct {
	// "insert 'true' if wish to include market goods under 'market goods' key or not. Such data can add a lot of extra weight"
	IncludeMarketGoods bool `protobuf:"varint,1,opt,name=include_market_goods,json=includeMarketGoods,proto3" json:"include_market_goods,omitempty"`
	// insert 'true' if wish to include tech compatibility data. can be adding a lot of extra weight
	IncludeTechCompat bool `protobuf:"varint,2,opt,name=include_tech_compat,json=includeTechCompat,proto3" json:"include_tech_compat,omitempty"`
	// Apply filtering same as darkstat does by default for its tab. Usually means showing only items that can be bought/crafted/or found
	FilterToUseful bool `protobuf:"varint,3,opt,name=filter_to_useful,json=filterToUseful,proto3" json:"filter_to_useful,omitempty"`
	// filters by item nicknames
	FilterNicknames []string `protobuf:"bytes,4,rep,name=filter_nicknames,json=filterNicknames,proto3" json:"filter_nicknames,omitempty"`
}

type GetCommoditiesInput struct {
	// To Include market goods, where the item is sold and bought or not. Adds a lot of extra weight to data
	//
	// Example: `false`
	IncludeMarketGoods bool `protobuf:"varint,1,opt,name=include_market_goods,json=includeMarketGoods,proto3" json:"include_market_goods,omitempty"`
	FilterToUseful     bool `protobuf:"varint,2,opt,name=filter_to_useful,json=filterToUseful,proto3" json:"filter_to_useful,omitempty"`
	// filters by item nicknames
	FilterNicknames []string `protobuf:"bytes,3,rep,name=filter_nicknames,json=filterNicknames,proto3" json:"filter_nicknames,omitempty"`
}

type GetBasesInput struct {
	// "insert 'true' if wish to include market goods under 'market goods' key or not. Such data can add a lot of extra weight"
	IncludeMarketGoods bool `protobuf:"varint,1,opt,name=include_market_goods,json=includeMarketGoods,proto3" json:"include_market_goods,omitempty"`
	// Apply filtering same as darkstat does by default for its tab. Usually means showing only items that can be bought/crafted/or found
	FilterToUseful bool `protobuf:"varint,2,opt,name=filter_to_useful,json=filterToUseful,proto3" json:"filter_to_useful,omitempty"`
	// filters by base nicknames
	FilterNicknames []string `protobuf:"bytes,3,rep,name=filter_nicknames,json=filterNicknames,proto3" json:"filter_nicknames,omitempty"`
	// filters market goods to specific category. valid categories are written in market goods in same named attribute.
	FilterMarketGoodCategory []string `protobuf:"bytes,4,rep,name=filter_market_good_category,json=filterMarketGoodCategory,proto3" json:"filter_market_good_category,omitempty"`
}

type Empty struct {
}

type NumString *string

type ShopItem struct {
	Nickname  string    `protobuf:"bytes,1,opt,name=nickname,proto3" json:"nickname"`
	Name      string    `protobuf:"bytes,2,opt,name=name,proto3" json:"name"`
	Category  string    `protobuf:"bytes,3,opt,name=category,proto3" json:"category"`
	Id        NumString `protobuf:"varint,4,opt,name=id,proto3" json:"id"`
	Quantity  NumString `protobuf:"varint,5,opt,name=quantity,proto3" json:"quantity"`
	Price     NumString `protobuf:"varint,6,opt,name=price,proto3" json:"price"`
	SellPrice NumString `protobuf:"varint,7,opt,name=sell_price,json=sellPrice,proto3" json:"sellPrice"`
	MinStock  NumString `protobuf:"varint,8,opt,name=min_stock,json=minStock,proto3" json:"minStock"`
	MaxStock  NumString `protobuf:"varint,9,opt,name=max_stock,json=maxStock,proto3" json:"maxStock"`
}

type PoBCore struct {
	Nickname       string    `protobuf:"bytes,1,opt,name=nickname,proto3" json:"nickname"`
	Name           string    `protobuf:"bytes,2,opt,name=name,proto3" json:"name"`
	Pos            *string   `protobuf:"bytes,3,opt,name=pos,proto3,oneof" json:"pos"`
	Level          NumString `protobuf:"varint,4,opt,name=level,proto3,oneof" json:"level"`
	Money          NumString `protobuf:"varint,5,opt,name=money,proto3,oneof" json:"money,omitempty"`
	Health         *float64  `protobuf:"fixed64,6,opt,name=health,proto3,oneof" json:"health"`
	DefenseMode    NumString `protobuf:"varint,7,opt,name=defense_mode,json=defenseMode,proto3,oneof" json:"defenseMode"`
	SystemNick     *string   `protobuf:"bytes,8,opt,name=system_nick,json=systemNick,proto3,oneof" json:"systemNick"`
	SystemName     *string   `protobuf:"bytes,9,opt,name=system_name,json=systemName,proto3,oneof" json:"systemName"`
	FactionNick    *string   `protobuf:"bytes,10,opt,name=faction_nick,json=factionNick,proto3,oneof" json:"factionNick"`
	FactionName    *string   `protobuf:"bytes,11,opt,name=faction_name,json=factionName,proto3,oneof" json:"factionName"`
	ForumThreadUrl *string   `protobuf:"bytes,12,opt,name=forum_thread_url,json=forumThreadUrl,proto3,oneof" json:"forumThreadUrl,omitempty"`
	CargoSpaceLeft NumString `protobuf:"varint,13,opt,name=cargo_space_left,json=cargoSpaceLeft,proto3,oneof" json:"cargoSpaceLeft"`
	BasePos        *Pos      `protobuf:"bytes,14,opt,name=base_pos,json=basePos,proto3,oneof" json:"basePos"`
	SectorCoord    *string   `protobuf:"bytes,15,opt,name=sector_coord,json=sectorCoord,proto3,oneof" json:"sectorCoord"`
	Region         *string   `protobuf:"bytes,16,opt,name=region,proto3,oneof" json:"region"`
}

type PoBGood struct {
	Nickname              string         `protobuf:"bytes,1,opt,name=nickname,proto3" json:"nickname"`
	Name                  string         `protobuf:"bytes,2,opt,name=name,proto3" json:"name"`
	TotalBuyableFromBases NumString      `protobuf:"varint,3,opt,name=total_buyable_from_bases,json=totalBuyableFromBases,proto3" json:"totalBuyableFromBases"`
	TotalSellableToBases  NumString      `protobuf:"varint,4,opt,name=total_sellable_to_bases,json=totalSellableToBases,proto3" json:"totalSellableToBases"`
	BestPriceToBuy        NumString      `protobuf:"varint,5,opt,name=best_price_to_buy,json=bestPriceToBuy,proto3,oneof" json:"bestPriceToBuy"`
	BestPriceToSell       NumString      `protobuf:"varint,6,opt,name=best_price_to_sell,json=bestPriceToSell,proto3,oneof" json:"bestPriceToSell"`
	Category              string         `protobuf:"bytes,7,opt,name=category,proto3" json:"category"`
	AnyBaseSells          bool           `protobuf:"varint,8,opt,name=any_base_sells,json=anyBaseSells,proto3" json:"anyBaseSells"`
	AnyBaseBuys           bool           `protobuf:"varint,9,opt,name=any_base_buys,json=anyBaseBuys,proto3" json:"anyBaseBuys"`
	Bases                 []*PoBGoodBase `protobuf:"bytes,10,rep,name=bases,proto3" json:"bases"`
	Volume                float64        `protobuf:"fixed64,11,opt,name=volume,proto3" json:"volume"`
	ShipClass             NumString      `protobuf:"varint,12,opt,name=ship_class,json=shipClass,proto3,oneof" json:"shipClass"`
}

type PoBGoodBase struct {
	ShopItem *ShopItem `protobuf:"bytes,1,opt,name=shop_item,json=shopItem,proto3" json:"shopItem"`
	Base     *PoBCore  `protobuf:"bytes,2,opt,name=base,proto3" json:"base"`
}

type GetPoBGoodsReply struct {
	Items []*PoBGood `protobuf:"bytes,1,rep,name=items,proto3" json:"items"`
}
