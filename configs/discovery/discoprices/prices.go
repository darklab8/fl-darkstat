package discoprices

import (
	"github.com/darklab8/fl-darkstat/configs/cfg"
	"github.com/darklab8/fl-darkstat/configs/configs_mapped/parserutils/filefind/file"
	"github.com/darklab8/fl-darkstat/configs/configs_mapped/parserutils/iniload"
	"github.com/darklab8/fl-darkstat/configs/configs_mapped/parserutils/semantic"
)

type Price struct {
	semantic.Model
	BaseNickname      *semantic.String
	CommodityNickname *semantic.String
	PriceBaseBuysFor  *semantic.Int
	PriceBaseSellsFor *semantic.Int
	BaseSells         *semantic.Bool
}

type Config struct {
	*iniload.IniLoader
	Prices       []*Price
	BasesPerGood map[string][]*Price
	GoodsPerBase map[string][]*Price
}

func Read(input_file *iniload.IniLoader) *Config {
	conf := &Config{
		IniLoader:    input_file,
		BasesPerGood: make(map[string][]*Price),
		GoodsPerBase: make(map[string][]*Price),
	}

	for _, price_info := range input_file.SectionMap["[price]"] {

		marketgood_key := cfg.Key("marketgood")
		for mg_index, _ := range price_info.ParamMap[marketgood_key] {
			market_good := &Price{}
			market_good.Map(price_info)
			market_good.BaseNickname = semantic.NewString(price_info, marketgood_key, semantic.WithLowercaseS(), semantic.WithoutSpacesS(), semantic.OptsS(semantic.Index(mg_index), semantic.Order(0)))
			market_good.CommodityNickname = semantic.NewString(price_info, marketgood_key, semantic.WithLowercaseS(), semantic.WithoutSpacesS(), semantic.OptsS(semantic.Index(mg_index), semantic.Order(1)))
			market_good.PriceBaseBuysFor = semantic.NewInt(price_info, marketgood_key, semantic.Index(mg_index), semantic.Order(2))
			market_good.PriceBaseSellsFor = semantic.NewInt(price_info, marketgood_key, semantic.Index(mg_index), semantic.Order(3))
			market_good.BaseSells = semantic.NewBool(price_info, marketgood_key, semantic.IntBool, semantic.Index(mg_index), semantic.Order(4))
			conf.Prices = append(conf.Prices, market_good)
			conf.BasesPerGood[market_good.CommodityNickname.Get()] = append(conf.BasesPerGood[market_good.CommodityNickname.Get()], market_good)
			conf.GoodsPerBase[market_good.BaseNickname.Get()] = append(conf.GoodsPerBase[market_good.BaseNickname.Get()], market_good)
		}
	}

	return conf
}

func (frelconfig *Config) Write() *file.File {
	return &file.File{}
}
