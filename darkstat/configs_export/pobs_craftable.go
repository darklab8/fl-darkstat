package configs_export

import (
	"strings"

	"github.com/darklab8/fl-darkstat/configs/cfg"
	"github.com/darklab8/fl-darkstat/configs/configs_mapped/parserutils/inireader"
	"github.com/darklab8/fl-darkstat/configs/discovery/base_recipe_items"
	"github.com/darklab8/fl-darkstat/darkstat/configs_export/infocarder"
	"github.com/darklab8/go-utils/utils/ptr"
)

func (e *Exporter) pob_produced() map[string]bool {
	if e.craftable_cached != nil {
		return e.craftable_cached
	}

	e.craftable_cached = make(map[string]bool)

	if e.Mapped.Discovery != nil {
		for _, recipe := range e.Mapped.Discovery.BaseRecipeItems.Recipes {
			for _, produced := range recipe.ProcucedItem {
				e.craftable_cached[produced.Get()] = true
			}
		}
	}

	if e.Mapped.FLSR != nil {
		for _, recipe := range e.Mapped.FLSR.FLSRRecipes.Products {
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
		Name:               e.Mapped.CraftableBaseName(),
		MarketGoodsPerNick: make(map[CommodityKey]*MarketGood),
		Nickname:           cfg.BaseUniNick(pob_crafts_nickname),
		SystemNickname:     "neverwhere",
		System:             "Neverwhere",
		Region:             "NEVERWHERE",
		FactionName:        "Player Crafts",
	}

	base.Archetypes = append(base.Archetypes, pob_crafts_nickname)

	for produced, _ := range pob_produced {
		market_good := &MarketGood{
			GoodInfo:             e.GetGoodInfo(produced),
			BaseSells:            true,
			IsServerSideOverride: true,
		}

		market_good_key := GetCommodityKey(market_good.Nickname, market_good.ShipClass)
		base.MarketGoodsPerNick[market_good_key] = market_good

		var infocard_addition infocarder.InfocardBuilder
		if e.Mapped.Discovery != nil {
			var any_recipe *base_recipe_items.CommodityRecipe
			if recipes, ok := e.Mapped.Discovery.BaseRecipeItems.RecipePerProduced[market_good.Nickname]; ok {

				infocard_addition.WriteLineStr(`CRAFTING RECIPES:`)
				for _, recipe := range recipes {
					sector := recipe.Model.RenderModel()
					infocard_addition.WriteLineStr(string(sector.OriginalType))
					for _, param := range sector.Params {
						infocard_addition.WriteLineStr(string(param.ToString(inireader.WithComments(false))))
					}
					infocard_addition.WriteLineStr("")
					any_recipe = recipe
				}
			}

			if any_recipe != nil {
				if craft_type, ok := any_recipe.CraftType.GetValue(); ok {
					if factory, ok := e.Mapped.Discovery.BaseRecipeModules.FactoryByCraftType[craft_type]; ok {
						module_name := factory.Infotext.Get()
						market_good.DiscoveryFactoryName = ptr.Ptr(module_name)
					}
				}
			}
		}
		if e.Mapped.FLSR != nil {
			if e.Mapped.FLSR.FLSRRecipes != nil {
				if recipes, ok := e.Mapped.FLSR.FLSRRecipes.ProductsByNick[market_good.Nickname]; ok {
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

		var info infocarder.InfocardBuilder
		if value, ok := e.GetInfocard2(infocarder.InfocardKey(market_good.Nickname)); ok {
			info.Lines = value
		}

		add_line_about_recipes := func(info infocarder.Infocard) infocarder.Infocard {
			add_line := func(index int, line infocarder.InfocardLine) {

				if len(info) > 0 && index < len(info) {
					info = append(info[:index+1], info[index:]...)
					info[index] = line
				} else {
					info = append(info, line)
				}
			}
			strip_line := func(line string) string {
				return strings.ReplaceAll(strings.ReplaceAll(line, " ", ""), "\u00a0", "")
			}
			if len(infocard_addition.Lines) > 0 {
				line_position := 1
				add_line(line_position, infocarder.InfocardLine{Phrases: []infocarder.InfocardPhrase{{Phrase: `Item has crafting recipes below`, Bold: true}}})
				if strip_line(info[0].ToStr()) != "" {
					add_line(1, infocarder.NewInfocardSimpleLine(""))
					line_position += 1
				}
				if strip_line(info[line_position+1].ToStr()) != "" {
					add_line(line_position+1, infocarder.NewInfocardSimpleLine(""))
				}
			}
			return info
		}
		info.Lines = add_line_about_recipes(info.Lines)

		e.PutInfocard(infocarder.InfocardKey(market_good.Nickname), append(info.Lines, infocard_addition.Lines...))

		if market_good.ShipNickname != "" {
			var info infocarder.Infocard
			if value, ok := e.GetInfocard2(infocarder.InfocardKey(market_good.ShipNickname)); ok {
				info = value
			}
			info = add_line_about_recipes(info)
			e.PutInfocard(infocarder.InfocardKey(market_good.ShipNickname), append(info, infocard_addition.Lines...))
		}
	}

	var sb infocarder.InfocardBuilder
	sb.WriteLineStr(base.Name)
	sb.WriteLineStr(`This is only pseudo base to show availability of player crafts`)
	sb.WriteLineStr(``)
	sb.WriteLineStr(`At the bottom of each item infocard it shows CRAFTING RECIPES`)

	e.PutInfocard(infocarder.InfocardKey(base.Nickname), sb.Lines)

	bases = append(bases, base)
	return bases
}
