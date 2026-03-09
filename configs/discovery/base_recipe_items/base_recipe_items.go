package base_recipe_items

import (
	"github.com/darklab8/fl-darkstat/configs/cfg"
	"github.com/darklab8/fl-darkstat/configs/configs_mapped/parserutils/filefind/file"
	"github.com/darklab8/fl-darkstat/configs/configs_mapped/parserutils/iniload"
	"github.com/darklab8/fl-darkstat/configs/configs_mapped/parserutils/semantic"
)

type AffiliationBonus struct {
	FactionNickname *semantic.String
	BonusMultiplier *semantic.Float
}

type ConsumedAlt struct {
	Amount *semantic.Int
	Items  []*semantic.String
}

type ProducedItem struct {
	Nickname *semantic.String
	Amount   *semantic.Int
}

type Produced struct {
	ProducedItem
	FactionProduced map[string]*ProducedItem
}

type Catalyst struct {
	Nickname *semantic.String
	Amount   *semantic.Int
}
type Consumed struct {
	Nickname *semantic.String
	Amount   *semantic.Int
}

type CommodityRecipe struct {
	semantic.Model
	Nickname     *semantic.String
	CraftType    *semantic.String
	ProducedItem []*Produced
	ConsumedItem []*Consumed

	Catalysts        []*Catalyst
	ConsumedAlt      []*ConsumedAlt
	ShortCutNum      *semantic.Int
	CookingRate      *semantic.Int
	RequiredLevel    *semantic.Int
	AffiliationBonus []*AffiliationBonus
	LoopProduction   *semantic.Bool
	Restricted       *semantic.Bool
}

type Config struct {
	*iniload.IniLoader
	Recipes           []*CommodityRecipe
	RecipePerConsumed map[string][]*CommodityRecipe
	RecipePerProduced map[string][]*CommodityRecipe
}

func Read(input_file *iniload.IniLoader) *Config {
	conf := &Config{
		IniLoader:         input_file,
		RecipePerConsumed: make(map[string][]*CommodityRecipe),
		RecipePerProduced: make(map[string][]*CommodityRecipe),
	}

	for _, recipe_info := range input_file.SectionMap["[recipe]"] {

		recipe := &CommodityRecipe{
			Nickname:       semantic.NewString(recipe_info, cfg.Key("nickname"), semantic.WithLowercaseS(), semantic.WithoutSpacesS()),
			CraftType:      semantic.NewString(recipe_info, cfg.Key("craft_type"), semantic.WithLowercaseS(), semantic.WithoutSpacesS()),
			CookingRate:    semantic.NewInt(recipe_info, cfg.Key("cooking_rate")),
			RequiredLevel:  semantic.NewInt(recipe_info, cfg.Key("reqlevel")),
			LoopProduction: semantic.NewBool(recipe_info, cfg.Key("loop_production"), semantic.IntBool),
			Restricted:     semantic.NewBool(recipe_info, cfg.Key("restricted"), semantic.StrBool),
			ShortCutNum:    semantic.NewInt(recipe_info, cfg.Key("shortcut_number")),
		}
		recipe.Map(recipe_info)

		for index, _ := range recipe_info.ParamMap[cfg.Key("produced_item")] {
			item := &Produced{
				ProducedItem: ProducedItem{
					Nickname: semantic.NewString(recipe_info, cfg.Key("produced_item"), semantic.WithLowercaseS(), semantic.WithoutSpacesS(), semantic.OptsS(semantic.Index(index))),
					Amount:   semantic.NewInt(recipe_info, cfg.Key("produced_item"), semantic.Order(1), semantic.Index(index)),
				},
			}
			recipe.ProducedItem = append(recipe.ProducedItem, item)
		}
		for index, _ := range recipe_info.ParamMap[cfg.Key("catalyst")] {
			item := &Catalyst{
				Nickname: semantic.NewString(recipe_info, cfg.Key("catalyst"), semantic.WithLowercaseS(), semantic.WithoutSpacesS(), semantic.OptsS(semantic.Index(index))),
				Amount:   semantic.NewInt(recipe_info, cfg.Key("catalyst"), semantic.Order(1), semantic.Index(index)),
			}
			recipe.Catalysts = append(recipe.Catalysts, item)
		}
		for produced_index, produced_affiliation_info := range recipe_info.ParamMap[cfg.Key("produced_affiliation")] {
			produced := &Produced{
				ProducedItem: ProducedItem{
					Nickname: semantic.NewString(recipe_info, cfg.Key("produced_affiliation"), semantic.WithLowercaseS(), semantic.WithoutSpacesS(), semantic.OptsS(semantic.Index(produced_index), semantic.Order(0))),
					Amount:   semantic.NewInt(recipe_info, cfg.Key("produced_affiliation"), semantic.Index(produced_index), semantic.Order(1)),
				},

				FactionProduced: make(map[string]*ProducedItem),
			}
			for i := 2; i < len(produced_affiliation_info.Values); i += 3 {
				faction := semantic.NewString(recipe_info, cfg.Key("produced_affiliation"), semantic.WithLowercaseS(), semantic.WithoutSpacesS(), semantic.OptsS(semantic.Index(produced_index), semantic.Order(i)))
				produced_instead := &ProducedItem{
					Nickname: semantic.NewString(recipe_info, cfg.Key("produced_affiliation"), semantic.WithLowercaseS(), semantic.WithoutSpacesS(), semantic.OptsS(semantic.Index(produced_index), semantic.Order(i+1))),
					Amount:   semantic.NewInt(recipe_info, cfg.Key("produced_affiliation"), semantic.Index(produced_index), semantic.Order(i+2)),
				}
				produced.FactionProduced[faction.Get()] = produced_instead
			}
			recipe.ProducedItem = append(recipe.ProducedItem, produced)
		}
		for index, _ := range recipe_info.ParamMap[cfg.Key("affiliation_bonus")] {
			bonus := &AffiliationBonus{
				FactionNickname: semantic.NewString(recipe_info, cfg.Key("affiliation_bonus"), semantic.WithLowercaseS(), semantic.WithoutSpacesS(), semantic.OptsS(semantic.Index(index))),
				BonusMultiplier: semantic.NewFloat(recipe_info, cfg.Key("affiliation_bonus"), semantic.Precision(2), semantic.OptsF(semantic.Order(1), semantic.Index(index))),
			}
			recipe.AffiliationBonus = append(recipe.AffiliationBonus, bonus)
		}

		for index, param := range recipe_info.ParamMap[cfg.Key("consumed_dynamic_alt")] {
			consum_alt := &ConsumedAlt{
				Amount: semantic.NewInt(recipe_info, cfg.Key("consumed_dynamic_alt"), semantic.Index(index)),
			}

			for order, _ := range param.Values {
				if order == 0 {
					continue
				}

				consum_alt.Items = append(consum_alt.Items,
					semantic.NewString(recipe_info,
						cfg.Key("consumed_dynamic_alt"),
						semantic.WithLowercaseS(),
						semantic.WithoutSpacesS(),
						semantic.OptsS(semantic.Index(index), semantic.Order(order)),
					),
				)
			}

			recipe.ConsumedAlt = append(recipe.ConsumedAlt, consum_alt)
		}

		for index, _ := range recipe_info.ParamMap[cfg.Key("consumed")] {
			item := &Consumed{
				Nickname: semantic.NewString(recipe_info, cfg.Key("consumed"), semantic.WithLowercaseS(), semantic.WithoutSpacesS(), semantic.OptsS(semantic.Index(index))),
				Amount:   semantic.NewInt(recipe_info, cfg.Key("consumed"), semantic.Order(1), semantic.Index(index)),
			}
			recipe.ConsumedItem = append(recipe.ConsumedItem, item)
		}

		conf.Recipes = append(conf.Recipes, recipe)
		for _, consumed := range recipe.ConsumedItem {
			conf.RecipePerConsumed[consumed.Nickname.Get()] = append(conf.RecipePerConsumed[consumed.Nickname.Get()], recipe)
		}
		for _, produced := range recipe.ProducedItem {
			conf.RecipePerProduced[produced.Nickname.Get()] = append(conf.RecipePerProduced[produced.Nickname.Get()], recipe)
		}
	}

	return conf
}

func (frelconfig *Config) Write() *file.File {
	return &file.File{}
}
