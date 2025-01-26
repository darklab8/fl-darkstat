package market_mapped

import (
	"github.com/darklab8/fl-darkstat/configs/cfg"
	"github.com/darklab8/fl-darkstat/configs/configs_mapped/parserutils/filefind/file"
	"github.com/darklab8/fl-darkstat/configs/configs_mapped/parserutils/iniload"
	"github.com/darklab8/fl-darkstat/configs/configs_mapped/parserutils/semantic"

	"github.com/darklab8/go-utils/utils/utils_types"
)

// Not implemented. Create SemanticMultiKeyValue
type MarketGood struct {
	semantic.Model
	Nickname *semantic.String // 0

	LevelRequired                       *semantic.Int   // 1
	RepRequired                         *semantic.Float // 2
	BaseSellsIPositiveAndDiscoSellPrice *semantic.Int   // 3

	baseSellsIfAboveZero *semantic.Int   // 4 . Base sells if value is above 0
	Display              *semantic.Int   // 5.
	PriceModifier        *semantic.Float // 6
}

func (m *MarketGood) BaseSells() bool {
	return m.BaseSellsIPositiveAndDiscoSellPrice.Get() > 0 && m.baseSellsIfAboveZero.Get() > 0
}

type BaseGood struct {
	semantic.Model
	Base *semantic.String

	MarketGoods    []*MarketGood
	MarketGoodsMap map[string]*MarketGood
}

type MarketGoodAtBase struct {
	MarketGood *MarketGood
	Base       cfg.BaseUniNick
}

type Config struct {
	Files []*iniload.IniLoader

	BaseGoods    []*BaseGood
	BasesPerGood map[string][]*MarketGoodAtBase
	GoodsPerBase map[cfg.BaseUniNick]*BaseGood
}

const (
	FILENAME_SHIPS       utils_types.FilePath = "market_ships.ini"
	FILENAME_COMMODITIES utils_types.FilePath = "market_commodities.ini"
	FILENAME_MISC        utils_types.FilePath = "market_misc.ini"
	BaseGoodType                              = "[basegood]"
)

var (
	KEY_MISSMATCH_SYSTEM_FILE = cfg.Key("missmatched_universe_system_and_file")
	KEY_MARKET_GOOD           = cfg.Key("marketgood")
	KEY_BASE                  = cfg.Key("base")
)

func Read(files []*iniload.IniLoader) *Config {
	frelconfig := &Config{
		Files:        files,
		GoodsPerBase: make(map[cfg.BaseUniNick]*BaseGood),
	}
	frelconfig.BaseGoods = make([]*BaseGood, 0)
	frelconfig.BasesPerGood = make(map[string][]*MarketGoodAtBase)

	for _, file := range frelconfig.Files {

		for _, section := range file.Sections {
			base_to_add := &BaseGood{
				MarketGoodsMap: make(map[string]*MarketGood),
			}
			base_to_add.Map(section)
			base_to_add.Base = semantic.NewString(section, KEY_BASE, semantic.WithLowercaseS(), semantic.WithoutSpacesS())
			base_nickname := cfg.BaseUniNick(base_to_add.Base.Get())

			for good_index, _ := range section.ParamMap[KEY_MARKET_GOOD] {
				good_to_add := &MarketGood{}
				good_to_add.Map(section)
				good_to_add.Nickname = semantic.NewString(section, KEY_MARKET_GOOD, semantic.OptsS(semantic.Index(good_index)), semantic.WithLowercaseS(), semantic.WithoutSpacesS())
				good_to_add.LevelRequired = semantic.NewInt(section, KEY_MARKET_GOOD, semantic.Index(good_index), semantic.Order(1))
				good_to_add.RepRequired = semantic.NewFloat(section, KEY_MARKET_GOOD, semantic.Precision(2), semantic.OptsF(semantic.Index(good_index), semantic.Order(2)))
				good_to_add.BaseSellsIPositiveAndDiscoSellPrice = semantic.NewInt(section, KEY_MARKET_GOOD, semantic.Index(good_index), semantic.Order(3))
				good_to_add.baseSellsIfAboveZero = semantic.NewInt(section, KEY_MARKET_GOOD, semantic.Index(good_index), semantic.Order(4))

				good_to_add.PriceModifier = semantic.NewFloat(section, KEY_MARKET_GOOD, semantic.Precision(2), semantic.OptsF(semantic.Index(good_index), semantic.Order(6)))
				base_to_add.MarketGoods = append(base_to_add.MarketGoods, good_to_add)
				base_to_add.MarketGoodsMap[good_to_add.Nickname.Get()] = good_to_add

				frelconfig.BasesPerGood[good_to_add.Nickname.Get()] = append(frelconfig.BasesPerGood[good_to_add.Nickname.Get()], &MarketGoodAtBase{
					MarketGood: good_to_add,
					Base:       base_nickname,
				})
			}

			frelconfig.BaseGoods = append(frelconfig.BaseGoods, base_to_add)
			frelconfig.GoodsPerBase[base_nickname] = base_to_add
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
