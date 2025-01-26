package base_recipe_items

import (
	"github.com/darklab8/fl-darkstat/configs/cfg"
	"github.com/darklab8/fl-darkstat/configs/configs_mapped/parserutils/filefind/file"
	"github.com/darklab8/fl-darkstat/configs/configs_mapped/parserutils/iniload"
	"github.com/darklab8/fl-darkstat/configs/configs_mapped/parserutils/semantic"
)

type CommodityRecipe struct {
	semantic.Model
	Nickname     *semantic.String
	ProcucedItem []*semantic.String
	ConsumedItem []*semantic.String
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
			Nickname: semantic.NewString(recipe_info, cfg.Key("nickname"), semantic.WithLowercaseS(), semantic.WithoutSpacesS()),
		}
		recipe.Map(recipe_info)

		for produced_index, _ := range recipe_info.ParamMap[cfg.Key("produced_item")] {

			recipe.ProcucedItem = append(recipe.ProcucedItem,
				semantic.NewString(recipe_info, cfg.Key("produced_item"), semantic.WithLowercaseS(), semantic.WithoutSpacesS(), semantic.OptsS(semantic.Index(produced_index))))
		}
		for produced_index, produced_affiliation_info := range recipe_info.ParamMap[cfg.Key("produced_affiliation")] {
			for i := 0; i < len(produced_affiliation_info.Values); i += 3 {
				recipe.ProcucedItem = append(recipe.ProcucedItem,
					semantic.NewString(recipe_info, cfg.Key("produced_affiliation"), semantic.WithLowercaseS(), semantic.WithoutSpacesS(), semantic.OptsS(semantic.Index(produced_index), semantic.Order(i))))
			}
		}

		for consumed_index, _ := range recipe_info.ParamMap[cfg.Key("consumed")] {

			recipe.ConsumedItem = append(recipe.ConsumedItem,
				semantic.NewString(recipe_info, cfg.Key("consumed"), semantic.WithLowercaseS(), semantic.WithoutSpacesS(), semantic.OptsS(semantic.Index(consumed_index))))

		}
		conf.Recipes = append(conf.Recipes, recipe)
		for _, consumed := range recipe.ConsumedItem {
			conf.RecipePerConsumed[consumed.Get()] = append(conf.RecipePerConsumed[consumed.Get()], recipe)
		}
		for _, produced := range recipe.ProcucedItem {
			conf.RecipePerProduced[produced.Get()] = append(conf.RecipePerProduced[produced.Get()], recipe)
		}
	}

	return conf
}

func (frelconfig *Config) Write() *file.File {
	return &file.File{}
}
