package base_recipe_modules

import (
	"github.com/darklab8/fl-darkstat/configs/cfg"
	"github.com/darklab8/fl-darkstat/configs/configs_mapped/parserutils/filefind/file"
	"github.com/darklab8/fl-darkstat/configs/configs_mapped/parserutils/iniload"
	"github.com/darklab8/fl-darkstat/configs/configs_mapped/parserutils/semantic"
)

type CommodityRecipe struct {
	semantic.Model
	Nickname  *semantic.String
	CraftList []*semantic.String
	Infotext  *semantic.String
}

type Config struct {
	*iniload.IniLoader
	Recipes            []*CommodityRecipe
	FactoryByCraftType map[string]*CommodityRecipe
}

func Read(input_file *iniload.IniLoader) *Config {
	conf := &Config{
		IniLoader:          input_file,
		FactoryByCraftType: make(map[string]*CommodityRecipe),
	}

	for _, recipe_info := range input_file.SectionMap["[recipe]"] {

		recipe := &CommodityRecipe{
			Nickname: semantic.NewString(recipe_info, cfg.Key("nickname"), semantic.WithLowercaseS(), semantic.WithoutSpacesS()),
			Infotext: semantic.NewString(recipe_info, cfg.Key("infotext")),
		}

		recipe.Map(recipe_info)

		for craft_list_index, _ := range recipe_info.ParamMap[cfg.Key("craft_list")] {

			recipe.CraftList = append(recipe.CraftList,
				semantic.NewString(recipe_info, cfg.Key("craft_list"), semantic.WithLowercaseS(), semantic.WithoutSpacesS(), semantic.OptsS(semantic.Index(craft_list_index))))
		}

		conf.Recipes = append(conf.Recipes, recipe)
		for _, craft_type := range recipe.CraftList {
			conf.FactoryByCraftType[craft_type.Get()] = recipe
		}
	}

	return conf
}

func (frelconfig *Config) Write() *file.File {
	return &file.File{}
}
