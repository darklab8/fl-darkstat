package equipment_mapped

import (
	"github.com/darklab8/fl-darkstat/configs/configs_mapped/parserutils/filefind/file"
	"github.com/darklab8/fl-darkstat/configs/configs_mapped/parserutils/iniload"
	"github.com/darklab8/fl-darkstat/configs/configs_mapped/parserutils/semantic"

	"github.com/darklab8/fl-darkstat/configs/cfg"
	"github.com/darklab8/go-utils/utils/utils_types"
)

type Commodity struct {
	semantic.Model
	Nickname  *semantic.String
	Equipment *semantic.String
	Category  *semantic.String

	Price         *semantic.Int
	Combinable    *semantic.Bool
	GoodSellPrice *semantic.Float
	BadBuyPrice   *semantic.Float
	BadSellPrice  *semantic.Float
	GoodBuyPrice  *semantic.Float
	ShopArchetype *semantic.Path
	ItemIcon      *semantic.Path
	JumpDist      *semantic.Int
}

type Ship struct {
	semantic.Model
	Category *semantic.String
	Nickname *semantic.String
	Hull     *semantic.String
	Addons   []*Addon
}

type Addon struct {
	semantic.Model
	ItemNickname *semantic.String
	ItemClass    *semantic.String
	Quantity     *semantic.Int
}
type ShipHull struct {
	semantic.Model
	Nickname *semantic.String
	Category *semantic.String
	Ship     *semantic.String
	Price    *semantic.Int
	IdsName  *semantic.Int
}

type Good struct {
	semantic.Model
	Category *semantic.String
	Nickname *semantic.String
	Price    *semantic.Int
}

type Config struct {
	Files []*iniload.IniLoader

	Goods    []*Good
	GoodsMap map[string]*Good

	Commodities        []*Commodity
	CommoditiesMap     map[string]*Commodity
	Ships              []*Ship
	ShipsMap           map[string]*Ship
	ShipsMapByHull     map[string][]*Ship
	ShipHulls          []*ShipHull
	ShipHullsMap       map[string]*ShipHull
	ShipHullsMapByShip map[string]*ShipHull
}

const (
	FILENAME utils_types.FilePath = "goods.ini"
)

func Read(configs []*iniload.IniLoader) *Config {
	frelconfig := &Config{Files: configs}
	frelconfig.Commodities = make([]*Commodity, 0, 100)
	frelconfig.CommoditiesMap = make(map[string]*Commodity)
	frelconfig.Ships = make([]*Ship, 0, 100)
	frelconfig.ShipsMap = make(map[string]*Ship)
	frelconfig.ShipHulls = make([]*ShipHull, 0, 100)
	frelconfig.ShipHullsMap = make(map[string]*ShipHull)

	frelconfig.Goods = make([]*Good, 0, 100)
	frelconfig.GoodsMap = make(map[string]*Good)
	frelconfig.ShipHullsMapByShip = make(map[string]*ShipHull)
	frelconfig.ShipsMapByHull = make(map[string][]*Ship)

	for _, config := range configs {
		for _, section := range config.SectionMap["[good]"] {
			good := &Good{}
			good.Map(section)
			good.Nickname = semantic.NewString(section, cfg.Key("nickname"), semantic.WithLowercaseS(), semantic.WithoutSpacesS())
			good.Category = semantic.NewString(section, cfg.Key("category"), semantic.WithLowercaseS(), semantic.WithoutSpacesS())
			good.Price = semantic.NewInt(section, cfg.Key("price"), semantic.Optional())
			frelconfig.Goods = append(frelconfig.Goods, good)
			frelconfig.GoodsMap[good.Nickname.Get()] = good

			category := good.Category.Get()
			switch category {
			case "commodity":
				commodity := &Commodity{}
				commodity.Map(section)
				commodity.Category = semantic.NewString(section, cfg.Key("category"), semantic.WithLowercaseS(), semantic.WithoutSpacesS())
				commodity.Nickname = semantic.NewString(section, cfg.Key("nickname"), semantic.WithLowercaseS(), semantic.WithoutSpacesS())
				commodity.Equipment = semantic.NewString(section, cfg.Key("equipment"), semantic.WithLowercaseS(), semantic.WithoutSpacesS())
				commodity.Price = semantic.NewInt(section, cfg.Key("price"))
				commodity.Combinable = semantic.NewBool(section, cfg.Key("combinable"), semantic.StrBool)
				commodity.GoodSellPrice = semantic.NewFloat(section, cfg.Key("good_sell_price"), semantic.Precision(2))
				commodity.BadBuyPrice = semantic.NewFloat(section, cfg.Key("bad_buy_price"), semantic.Precision(2))
				commodity.BadSellPrice = semantic.NewFloat(section, cfg.Key("bad_sell_price"), semantic.Precision(2))
				commodity.GoodBuyPrice = semantic.NewFloat(section, cfg.Key("good_buy_price"), semantic.Precision(2))
				commodity.ShopArchetype = semantic.NewPath(section, cfg.Key("shop_archetype"))
				commodity.ItemIcon = semantic.NewPath(section, cfg.Key("item_icon"))
				commodity.JumpDist = semantic.NewInt(section, cfg.Key("jump_dist"))

				frelconfig.Commodities = append(frelconfig.Commodities, commodity)
				frelconfig.CommoditiesMap[commodity.Nickname.Get()] = commodity
			case "ship":
				ship := &Ship{}
				ship.Map(section)
				ship.Category = semantic.NewString(section, cfg.Key("category"))
				ship.Nickname = semantic.NewString(section, cfg.Key("nickname"), semantic.WithLowercaseS(), semantic.WithoutSpacesS())
				ship.Hull = semantic.NewString(section, cfg.Key("hull"))

				for addon_i, _ := range section.ParamMap[cfg.Key("addon")] {
					addon := &Addon{
						ItemNickname: semantic.NewString(section, cfg.Key("addon"),
							semantic.OptsS(semantic.Index(addon_i), semantic.Order(0)),
							semantic.WithLowercaseS(), semantic.WithoutSpacesS()),
						ItemClass: semantic.NewString(section, cfg.Key("addon"),
							semantic.OptsS(semantic.Index(addon_i), semantic.Order(1)),
							semantic.WithLowercaseS(), semantic.WithoutSpacesS()),
						Quantity: semantic.NewInt(section, cfg.Key("addon"),
							semantic.Index(addon_i), semantic.Order(2),
						),
					}
					addon.Map(section)
					ship.Addons = append(ship.Addons, addon)
				}

				frelconfig.Ships = append(frelconfig.Ships, ship)
				frelconfig.ShipsMap[ship.Nickname.Get()] = ship
				frelconfig.ShipsMapByHull[ship.Hull.Get()] = append(frelconfig.ShipsMapByHull[ship.Hull.Get()], ship)
			case "shiphull":
				shiphull := &ShipHull{}
				shiphull.Map(section)
				shiphull.Category = semantic.NewString(section, cfg.Key("category"), semantic.WithLowercaseS(), semantic.WithoutSpacesS())
				shiphull.Nickname = semantic.NewString(section, cfg.Key("nickname"), semantic.WithLowercaseS(), semantic.WithoutSpacesS())
				shiphull.Ship = semantic.NewString(section, cfg.Key("ship"), semantic.WithLowercaseS(), semantic.WithoutSpacesS())
				shiphull.Price = semantic.NewInt(section, cfg.Key("price"))
				shiphull.IdsName = semantic.NewInt(section, cfg.Key("ids_name"), semantic.Optional())

				frelconfig.ShipHulls = append(frelconfig.ShipHulls, shiphull)
				frelconfig.ShipHullsMap[shiphull.Nickname.Get()] = shiphull
				frelconfig.ShipHullsMapByShip[shiphull.Ship.Get()] = shiphull
			}

		}
	}

	return frelconfig
}

func (frelconfig *Config) Write() []*file.File {
	var files []*file.File
	for _, file := range frelconfig.Files {
		inifile := file.Render()
		inifile.Write(inifile.File)
		files = append(files, inifile.File)
	}
	return files
}
