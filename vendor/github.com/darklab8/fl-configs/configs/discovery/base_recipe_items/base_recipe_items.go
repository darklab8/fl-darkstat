package base_recipe_items

import (
	"github.com/darklab8/fl-configs/configs/configs_mapped/parserutils/filefind/file"
	"github.com/darklab8/fl-configs/configs/configs_mapped/parserutils/iniload"
	"github.com/darklab8/fl-configs/configs/configs_mapped/parserutils/semantic"
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
}

func Read(input_file *iniload.IniLoader) *Config {
	conf := &Config{
		IniLoader:         input_file,
		RecipePerConsumed: make(map[string][]*CommodityRecipe),
	}

	for _, recipe_info := range input_file.SectionMap["[recipe]"] {

		recipe := &CommodityRecipe{
			Nickname: semantic.NewString(recipe_info, "nickname", semantic.WithLowercaseS(), semantic.WithoutSpacesS()),
		}
		recipe.Map(recipe_info)

		for produced_index, _ := range recipe_info.ParamMap["produced_item"] {

			recipe.ProcucedItem = append(recipe.ProcucedItem,
				semantic.NewString(recipe_info, "produced_item", semantic.WithLowercaseS(), semantic.WithoutSpacesS(), semantic.OptsS(semantic.Index(produced_index))))
		}

		for consumed_index, _ := range recipe_info.ParamMap["consumed"] {

			recipe.ConsumedItem = append(recipe.ConsumedItem,
				semantic.NewString(recipe_info, "consumed", semantic.WithLowercaseS(), semantic.WithoutSpacesS(), semantic.OptsS(semantic.Index(consumed_index))))

		}
		conf.Recipes = append(conf.Recipes, recipe)
		for _, consumed := range recipe.ConsumedItem {
			conf.RecipePerConsumed[consumed.Get()] = append(conf.RecipePerConsumed[consumed.Get()], recipe)
		}

	}

	return conf
}

func (frelconfig *Config) Write() *file.File {
	return &file.File{}
}
