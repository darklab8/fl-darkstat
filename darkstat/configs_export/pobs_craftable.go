package configs_export

import (
	"strings"

	"github.com/darklab8/fl-darkstat/configs/cfg"
	"github.com/darklab8/fl-darkstat/configs/configs_mapped/parserutils/inireader"
)

func (e *Exporter) pob_produced() map[string]bool {
	if e.craftable_cached != nil {
		return e.craftable_cached
	}

	e.craftable_cached = make(map[string]bool)

	if e.Configs.Discovery != nil {
		for _, recipe := range e.Configs.Discovery.BaseRecipeItems.Recipes {
			for _, produced := range recipe.ProcucedItem {
				e.craftable_cached[produced.Get()] = true
			}
		}
	}

	if e.Configs.FLSR != nil {
		for _, recipe := range e.Configs.FLSR.FLSRRecipes.Products {
			e.craftable_cached[recipe.Product.Get()] = true
		}
	}

	return e.craftable_cached
}

const (
	pob_crafts_nickname = "crafts"
)

func (e *Exporter) EnhanceBasesWithPobCrafts(bases []*Base) []*Base {
	pob_produced := e.pob_produced()

	base := &Base{
		Name:               e.Configs.CraftableBaseName(),
		MarketGoodsPerNick: make(map[CommodityKey]*MarketGood),
		Nickname:           cfg.BaseUniNick(pob_crafts_nickname),
		InfocardKey:        InfocardKey(pob_crafts_nickname),
		SystemNickname:     "neverwhere",
		System:             "Neverwhere",
		Region:             "Neverwhere",
		FactionName:        "Player Crafts",
	}

	base.Archetypes = append(base.Archetypes, pob_crafts_nickname)

	for produced, _ := range pob_produced {
		market_good := &MarketGood{
			GoodInfo:             e.GetGoodInfo(produced),
			BaseSells:            true,
			ShipClass:            -1,
			IsServerSideOverride: true,
		}
		e.Hashes[market_good.Nickname] = market_good.NicknameHash

		market_good_key := GetCommodityKey(market_good.Nickname, market_good.ShipClass)
		base.MarketGoodsPerNick[market_good_key] = market_good

		var infocard_addition InfocardBuilder
		if e.Configs.Discovery != nil {
			if recipes, ok := e.Configs.Discovery.BaseRecipeItems.RecipePerProduced[market_good.Nickname]; ok {
				infocard_addition.WriteLineStr(`CRAFTING RECIPES:`)
				for _, recipe := range recipes {
					sector := recipe.Model.RenderModel()
					infocard_addition.WriteLineStr(string(sector.OriginalType))
					for _, param := range sector.Params {
						infocard_addition.WriteLineStr(string(param.ToString(inireader.WithComments(false))))
					}
					infocard_addition.WriteLineStr("")
				}
			}
		}
		if e.Configs.FLSR != nil {
			if e.Configs.FLSR.FLSRRecipes != nil {
				if recipes, ok := e.Configs.FLSR.FLSRRecipes.ProductsByNick[market_good.Nickname]; ok {
					infocard_addition.WriteLineStr(`CRAFTING RECIPES:`)
					for _, recipe := range recipes {
						sector := recipe.Model.RenderModel()
						infocard_addition.WriteLineStr(string(sector.OriginalType))
						for _, param := range sector.Params {
							infocard_addition.WriteLineStr(string(param.ToString(inireader.WithComments(false))))
						}
						infocard_addition.WriteLineStr("")
					}
				}
			}
		}

		var info InfocardBuilder
		if value, ok := e.Infocards[InfocardKey(market_good.Nickname)]; ok {
			info.Lines = value
		}

		add_line_about_recipes := func(info Infocard) Infocard {
			add_line := func(index int, line InfocardLine) {
				info = append(info[:index+1], info[index:]...)
				info[index] = line
			}
			strip_line := func(line string) string {
				return strings.ReplaceAll(strings.ReplaceAll(line, " ", ""), "\u00a0", "")
			}
			if len(infocard_addition.Lines) > 0 {
				line_position := 1
				add_line(line_position, InfocardLine{Phrases: []InfocardPhrase{{Phrase: `Item has crafting recipes below`, Bold: true}}})
				if strip_line(info[0].ToStr()) != "" {
					add_line(1, NewInfocardSimpleLine(""))
					line_position += 1
				}
				if strip_line(info[line_position+1].ToStr()) != "" {
					add_line(line_position+1, NewInfocardSimpleLine(""))
				}
			}
			return info
		}
		info.Lines = add_line_about_recipes(info.Lines)

		e.Infocards[InfocardKey(market_good.Nickname)] = append(info.Lines, infocard_addition.Lines...)

		if market_good.ShipNickname != "" {
			var info Infocard
			if value, ok := e.Infocards[InfocardKey(market_good.ShipNickname)]; ok {
				info = value
			}
			info = add_line_about_recipes(info)
			e.Infocards[InfocardKey(market_good.ShipNickname)] = append(info, infocard_addition.Lines...)
		}
	}

	var sb InfocardBuilder
	sb.WriteLineStr(base.Name)
	sb.WriteLineStr(`This is only pseudo base to show availability of player crafts`)
	sb.WriteLineStr(``)
	sb.WriteLineStr(`At the bottom of each item infocard it shows CRAFTING RECIPES`)

	e.Infocards[InfocardKey(base.Nickname)] = sb.Lines

	bases = append(bases, base)
	return bases
}
