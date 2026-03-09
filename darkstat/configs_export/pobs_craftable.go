package configs_export

import (
	"fmt"
	"math"
	"strings"

	"github.com/darklab8/fl-darkstat/configs/cfg"
	"github.com/darklab8/fl-darkstat/configs/configs_mapped/freelancer_mapped/data_mapped/universe_mapped"
	"github.com/darklab8/fl-darkstat/configs/configs_mapped/parserutils/inireader"
	"github.com/darklab8/fl-darkstat/configs/configs_settings/logus"
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
			for _, produced := range recipe.ProducedItem {
				e.craftable_cached[produced.Nickname.Get()] = true

				for _, alternates := range produced.FactionProduced {
					e.craftable_cached[alternates.Nickname.Get()] = true
				}
			}
		}
	}

	if e.Mapped.FLSR != nil {
		for _, recipe := range e.Mapped.FLSR.FLSRRecipes.Products {
			for _, produced := range recipe.Product {
				e.craftable_cached[produced.Nickname.Get()] = true
			}
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
			if recipes, ok := e.Mapped.Discovery.BaseRecipeItems.RecipePerProduced[market_good.Nickname]; ok {
				for _, recipe := range recipes {
					craft := CraftableDiscoInfo{}

					if command_group, ok := recipe.CraftType.GetValue(); ok {
						if command_shortcut, ok2 := recipe.ShortCutNum.GetValue(); ok2 {
							craft.Command = string(fmt.Sprintf("/craft %s %d", command_group, command_shortcut))
						}
					}

					if craft_type, ok := recipe.CraftType.GetValue(); ok {
						if factory, ok := e.Mapped.Discovery.BaseRecipeModules.FactoryByCraftType[craft_type]; ok {
							module_name := factory.Infotext.Get()
							craft.FactoryName = ptr.Ptr(module_name)
						}
					}

					if required_level, ok := recipe.RequiredLevel.GetValue(); ok {
						craft.RequiredLevel = ptr.Ptr(required_level)
					}

					if does_loop, ok := recipe.LoopProduction.GetValue(); ok {
						craft.LoopProduction = does_loop
					}

					{
						amount_volume_to_cook := 0.0
						for _, item := range recipe.ConsumedItem {
							name := item.Nickname.Get()
							equip := e.Mapped.Equip().ItemsMap[name]
							amount_volume_to_cook += equip.Volume.Get() * float64(item.Amount.Get())
						}
						for _, item := range recipe.ConsumedAlt {
							item_volume := 0.0
							for _, nickname := range item.Items {
								name := nickname.Get()
								equip := e.Mapped.Equip().ItemsMap[name]

								item_volume += equip.Volume.Get() * float64(item.Amount.Get())
							}
							item_volume = item_volume / float64(len(item.Items))
							amount_volume_to_cook += item_volume
						}

						cooking_rate := recipe.CookingRate.Get()
						craft.TotalVolume = amount_volume_to_cook
						craft.CookMinutes = amount_volume_to_cook / float64(cooking_rate)
					}

					market_good.CraftableInfos = append(market_good.CraftableInfos, CraftableInfo{Disco: craft})
				}
			}
			if recipes, ok := e.Mapped.Discovery.BaseRecipeItems.RecipePerProduced[market_good.Nickname]; ok {

				infocard_addition.WriteLineStr(`CRAFTING RECIPES:`)
				for index, recipe := range recipes {
					craft := market_good.CraftableInfos[index].Disco
					infocard_addition.WriteLineStrBold(string(fmt.Sprintf("[RECIPE #%d] (translated)", index+1)))

					if craft.Command != "" {
						infocard_addition.WriteLineStr(string(fmt.Sprintf("command: %s", craft.Command)))
					}
					if craft.FactoryName != nil {
						infocard_addition.WriteLineStr(string(fmt.Sprintf("Factory: %s", *craft.FactoryName)))
					}

					cooking_rate := recipe.CookingRate.Get()
					infocard_addition.WriteLineStr(string(fmt.Sprintf("Cooking: %d volume in minute", cooking_rate)))

					var sb_time strings.Builder
					sb_time.WriteString(fmt.Sprintf("Total recipe time: %.0f minutes", craft.CookMinutes))
					if math.Floor(craft.CookMinutes/60) > 0 {
						sb_time.WriteString(" [")
						sb_time.WriteString(fmt.Sprintf("%2.0fh - ", math.Floor(craft.CookMinutes/60)))
						sb_time.WriteString(fmt.Sprintf("%2.0fm", craft.CookMinutes-60*math.Floor(craft.CookMinutes/60)))
						sb_time.WriteString("]")
					}
					infocard_addition.WriteLineStr(sb_time.String())

					infocard_addition.WriteLineStr(string(fmt.Sprintf("Total recipe volume: %.0f", craft.TotalVolume)))

					if level, ok := recipe.RequiredLevel.GetValue(); ok {
						infocard_addition.WriteLineStr(string(fmt.Sprintf("Required core level: %d", level)))
					}
					if does_loop, ok := recipe.LoopProduction.GetValue(); ok {
						if does_loop {
							infocard_addition.WriteLineStr(string("recipe does loop in production"))
						}
					}

					for _, catalyst := range recipe.Catalysts {
						name := catalyst.Nickname.Get()
						if equip, ok := e.Mapped.Equip().ItemsMap[name]; ok {
							name = e.GetInfocardName(equip.IdsName.Get(), name)
						}
						amount := catalyst.Amount.Get()
						infocard_addition.WriteLineStr(string(fmt.Sprintf("* catalyst: %s (%d amount)", name, amount)))
					}
					for _, item := range recipe.ConsumedItem {
						name := item.Nickname.Get()
						volume := 0.0
						if equip, ok := e.Mapped.Equip().ItemsMap[name]; ok {
							name = e.GetInfocardName(equip.IdsName.Get(), name)
							volume = equip.Volume.Get()
						}
						amount := item.Amount.Get()
						infocard_addition.WriteLineStr(string(fmt.Sprintf("] consumed: %s (%d amount, %.0f vol)", name, amount, volume*float64(amount))))
					}
					for _, item := range recipe.ConsumedAlt {
						amount := item.Amount.Get()

						infocard_addition.WriteLineStr(string(fmt.Sprintf("] consumed one of next items (in %d amount):", amount)))

						for _, nickname := range item.Items {
							name := nickname.Get()
							volume := 0.0
							if equip, ok := e.Mapped.Equip().ItemsMap[name]; ok {
								name = e.GetInfocardName(equip.IdsName.Get(), name)
								volume = equip.Volume.Get()
							}
							infocard_addition.WriteLineStr(string(fmt.Sprintf("] --- %s (%.0f vol)", name, volume*float64(amount))))
						}
					}
					for _, item := range recipe.ProducedItem {
						name := item.Nickname.Get()
						if equip, ok := e.Mapped.Equip().ItemsMap[name]; ok {
							name = e.GetInfocardName(equip.IdsName.Get(), name)
						}
						amount := item.Amount.Get()
						var sb strings.Builder
						sb.WriteString(string(fmt.Sprintf("> produced: %s (%d amount)", name, amount)))
						if len(item.FactionProduced) > 0 {
							sb.WriteString(" or alternatively:")
						}

						infocard_addition.WriteLineStr(sb.String())

						for key, value := range item.FactionProduced {
							faction_name := key
							if group, ok := e.Mapped.InitialWorld.GroupsMap[faction_name]; ok {
								faction_name = e.GetInfocardName(group.IdsName.Get(), faction_name)
							}
							name := value.Nickname.Get()
							if equip, ok := e.Mapped.Equip().ItemsMap[name]; ok {
								name = e.GetInfocardName(equip.IdsName.Get(), name)
							}
							amount := value.Amount.Get()

							infocard_addition.WriteLineStr(string(fmt.Sprintf("> --- [%s (%d amount) if %s]", name, amount, faction_name)))
						}
					}

					if len(recipe.AffiliationBonus) > 0 {
						infocard_addition.WriteLineStr("# faction multipliers of needed mats:")
					}
					if restricted, ok := recipe.Restricted.GetValue(); ok {
						if restricted {
							infocard_addition.WriteLineStr("# !!! also restricted to only factions below:")
						}
					}
					for _, item := range recipe.AffiliationBonus {
						faction_name := item.FactionNickname.Get()
						if group, ok := e.Mapped.InitialWorld.GroupsMap[faction_name]; ok {
							faction_name = e.GetInfocardName(group.IdsName.Get(), faction_name)
						}
						amount := item.BonusMultiplier.Get()
						infocard_addition.WriteLineStr(string(fmt.Sprintf("# --- %s (%.2f multiplier)", faction_name, amount)))
					}

					infocard_addition.WriteLineStr("")

					infocard_addition.WriteLineStr(fmt.Sprintf("%s ; ([RECIPE #%d] original)", string(recipe.GetOriginalType()), index+1))
					for _, param := range GetModelWithoutLastComments(&recipe.Model) {
						infocard_addition.WriteLineStr(string(param.ToString(inireader.WithComments(false))))
					}
					infocard_addition.WriteLineStr("")
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
							} else {
								logus.Log.Error("craftable base has no name",
									typelog.Any("base_nickname", base_nickname),
								)
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

						for _, product := range recipe.Product {
							nickname := product.Nickname.Get()
							name := nickname
							if equip, ok := e.Mapped.Equip().ItemsMap[nickname]; ok {
								name = e.GetInfocardName(equip.IdsName.Get(), nickname)
							}
							recipe_info.Products = append(recipe_info.Products, Product{
								Name:     name,
								Amount:   product.Quantity.Get(),
								Nickname: nickname,
							})
						}

						command := strings.ReplaceAll(string(recipe.GetOriginalType()), "[", "")
						command = strings.ReplaceAll(command, "]", "")
						recipe_info.Command = fmt.Sprintf("/craft %s", command)

						market_good.CraftableInfos = append(market_good.CraftableInfos, CraftableInfo{FLSR: recipe_info})
					}

					infocard_addition.WriteLineStr(`CRAFTING RECIPES:`)
					for index, recipe := range market_good.CraftableInfos {
						infocard_addition.WriteLineStrBold(string(fmt.Sprintf("[RECIPE #%d] (translated)", index+1)))
						infocard_addition.WriteLineStr(string(fmt.Sprintf("command: %s", recipe.FLSR.Command)))
						infocard_addition.WriteLineStr(string(fmt.Sprintf("recipe cost: %d $ sirius credits", recipe.FLSR.CostPrice)))
						for _, ingredient := range recipe.FLSR.Ingredients {
							infocard_addition.WriteLineStr(string(fmt.Sprintf("ingredient: %s (%d amount)", ingredient.Name, ingredient.Amount)))
						}
						if len(recipe.FLSR.BaseNames) > 0 {
							for _, base := range recipe.FLSR.BaseNames {
								infocard_addition.WriteLineStr(string(fmt.Sprintf("base: %s", base)))
							}
						}

						for _, product := range recipe.FLSR.Products {
							infocard_addition.WriteLineStr(string(fmt.Sprintf("produces: %s (%d amount)", product.Name, product.Amount)))
						}

						infocard_addition.WriteLineStr("")

						original_recipe := recipes[index]
						infocard_addition.WriteLineStr(fmt.Sprintf("%s ; (#%d original)", string(original_recipe.GetOriginalType()), index+1))
						for _, param := range GetModelWithoutLastComments(&original_recipe.Model) {
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
