package configs_export

import (
	"fmt"
	"strings"

	"github.com/darklab8/fl-darkstat/configs/cfg"
	"github.com/darklab8/fl-darkstat/configs/configs_mapped/freelancer_mapped/data_mapped/universe_mapped"
	"github.com/darklab8/fl-darkstat/configs/configs_mapped/parserutils/inireader"
	"github.com/darklab8/fl-darkstat/configs/configs_settings/logus"
	"github.com/darklab8/fl-darkstat/configs/discovery/base_recipe_items"
	"github.com/darklab8/fl-darkstat/darkstat/configs_export/infocarder"
	"github.com/darklab8/go-utils/typelog"
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
		SystemNickname:     "readme",
		System:             "README",
		Region:             "README",
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
					infocard_addition.WriteLineStr(string(recipe.Model.GetOriginalType()))
					for _, param := range GetModelWithoutLastComments(&recipe.Model) {
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

					for _, recipe := range recipes {

						recipe_info := CraftableFLSRInfo{}

						for _, base_nickname := range recipe.BaseNicknames {
							base_nickname := base_nickname.Get()
							universe_base, ok := e.Mapped.Universe.BasesMap[universe_mapped.BaseNickname(base_nickname)]
							base_name := base_nickname
							if ok {
								base_name = e.GetInfocardName(universe_base.StridName.Get(), base_nickname)
							}
							recipe_info.BaseNames = append(recipe_info.BaseNames, base_name)
						}

						recipe_info.CostPrice, _ = recipe.Cost.GetValue()

						for _, ingredient := range recipe.Ingridients {
							nickname := ingredient.Nickname.Get()
							name := nickname
							if equip, ok := e.Mapped.Equip().ItemsMap[nickname]; ok {
								name = e.GetInfocardName(equip.IdsName.Get(), nickname)
							}
							recipe_info.Ingredients = append(recipe_info.Ingredients, Ingredient{
								Name:   name,
								Amount: ingredient.Quantity.Get(),
							})
						}

						command := strings.ReplaceAll(string(recipe.GetOriginalType()), "[", "")
						command = strings.ReplaceAll(command, "]", "")
						recipe_info.Command = command

						market_good.CraftableFLSRInfo = append(market_good.CraftableFLSRInfo, recipe_info)
					}

					infocard_addition.WriteLineStr(`CRAFTING RECIPES (translated):`)
					for index, recipe := range market_good.CraftableFLSRInfo {
						infocard_addition.WriteLineStr(string(fmt.Sprintf("[Recipe #%d]", index)))
						infocard_addition.WriteLineStr(string(fmt.Sprintf("command: /craft %s", recipe.Command)))
						infocard_addition.WriteLineStr(string(fmt.Sprintf("recipe cost: %d", recipe.CostPrice)))
						for _, ingredient := range recipe.Ingredients {
							infocard_addition.WriteLineStr(string(fmt.Sprintf("ingredient: %s (%d amount)", ingredient.Name, ingredient.Amount)))
						}
						if len(recipe.BaseNames) > 0 {
							for _, base := range recipe.BaseNames {
								infocard_addition.WriteLineStr(string(fmt.Sprintf("base: %s", base)))
							}
						}

						infocard_addition.WriteLineStr("")
					}

					infocard_addition.WriteLineStr(`CRAFTING RECIPES (original):`)
					for _, recipe := range recipes {
						infocard_addition.WriteLineStr(string(recipe.GetOriginalType()))
						for _, param := range GetModelWithoutLastComments(&recipe.Model) {
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
				defer func() {
					if err := recover(); err != nil {
						logus.Log.Error("badly added line", typelog.Any("err", err))
						info = append(info, line)
					}
				}()
				info = append(info[:index+1], info[index:]...)
				info[index] = line
			}
			strip_line := func(line string) string {
				return strings.ReplaceAll(strings.ReplaceAll(line, " ", ""), "\u00a0", "")
			}
			if len(infocard_addition.Lines) > 0 {
				line_position := 1
				add_line(line_position, infocarder.InfocardLine{Phrases: []infocarder.InfocardPhrase{{Phrase: `Item has CRAFTING RECIPES below`, Bold: true}}})
				if strip_line(info[0].ToStr()) != "" {
					add_line(1, infocarder.NewInfocardSimpleLine(""))
					line_position += 1
				}
				if line_position+1 < len(info) {
					if strip_line(info[line_position+1].ToStr()) != "" {
						add_line(line_position+1, infocarder.NewInfocardSimpleLine(""))
					}
				}

			}
			return info
		}
		info.Lines = add_line_about_recipes(info.Lines)

		TryPutBeforeIniConfigs := func(nickname infocarder.InfocardKey, original infocarder.Infocard, addition infocarder.Infocard) {

			beginning_of_ini_configs := -1
			for index, line := range original {
				str_line := line.ToStr()
				if strings.Contains(str_line, "[") && strings.Contains(str_line, "]") {
					beginning_of_ini_configs = index

					break
				}
			}

			if beginning_of_ini_configs != -1 {
				left := original[:beginning_of_ini_configs]  //store left slice elements in left variable
				right := original[beginning_of_ini_configs:] //store right slice elements in right variable

				addition_with_end := append(addition,
					infocarder.InfocardLine{Phrases: []infocarder.InfocardPhrase{{Phrase: `CRAFTING RECIPES FINISHED`}}},
					infocarder.InfocardLine{},
				)
				e.PutInfocard(infocarder.InfocardKey(nickname), append(left, append(addition_with_end, right...)...))
			} else {
				e.PutInfocard(infocarder.InfocardKey(nickname), append(original, addition...))
			}

		}

		TryPutBeforeIniConfigs(infocarder.InfocardKey(market_good.Nickname), info.Lines, infocard_addition.Lines)

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
	sb.WriteLineStr(``)
	sb.WriteLineStr(`Go to tab "BASES" and find it in the list there to see FULL LIST of possible all recipes!!!`)

	e.PutInfocard(infocarder.InfocardKey(base.Nickname), sb.Lines)

	e.CraftableBase = base
	bases = append(bases, base)
	return bases
}
