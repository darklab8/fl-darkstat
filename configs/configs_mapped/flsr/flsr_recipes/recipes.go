package flsr_recipes

import (
	"github.com/darklab8/fl-darkstat/configs/cfg"
	"github.com/darklab8/fl-darkstat/configs/configs_mapped/parserutils/filefind/file"
	"github.com/darklab8/fl-darkstat/configs/configs_mapped/parserutils/iniload"
	"github.com/darklab8/fl-darkstat/configs/configs_mapped/parserutils/semantic"
	"github.com/darklab8/go-utils/utils/utils_types"
)

type Ingredient struct {
	Nickname *semantic.String
	Quantity *semantic.Int
}

type Product struct {
	Nickname *semantic.String
	Quantity *semantic.Int
}

type Recipe struct {
	semantic.Model
	Product       []*Product
	Ingridients   []*Ingredient
	BaseNicknames []*semantic.String
	Cost          *semantic.Int
}

type Config struct {
	*iniload.IniLoader
	Products       []*Recipe
	ProductsByNick map[string][]*Recipe
}

const (
	FILENAME utils_types.FilePath = "flsr-crafting.cfg"
)

func Read(input_file *iniload.IniLoader) *Config {
	frelconfig := &Config{
		IniLoader:      input_file,
		ProductsByNick: make(map[string][]*Recipe),
	}

	for _, section := range input_file.Sections {

		recipe := &Recipe{
			Cost: semantic.NewInt(section, cfg.Key("cost")),
		}
		recipe.Map(section)

		if ingredients, ok := section.ParamMap["ingredient"]; ok {
			for index, _ := range ingredients {
				ingredient := &Ingredient{
					Nickname: semantic.NewString(section, cfg.Key("ingredient"), semantic.WithLowercaseS(), semantic.OptsS(semantic.Index(index)), semantic.WithoutSpacesS()),
					Quantity: semantic.NewInt(section, cfg.Key("ingredient"), semantic.Index(index), semantic.Order(1)),
				}
				recipe.Ingridients = append(recipe.Ingridients, ingredient)
			}
		}

		if products, ok := section.ParamMap["product"]; ok {
			for index, _ := range products {
				info := &Product{
					Nickname: semantic.NewString(section, cfg.Key("product"), semantic.WithLowercaseS(), semantic.OptsS(semantic.Index(index)), semantic.WithoutSpacesS()),
					Quantity: semantic.NewInt(section, cfg.Key("product"), semantic.Index(index), semantic.Order(1)),
				}
				recipe.Product = append(recipe.Product, info)
			}
		}

		if base_nicknames, ok := section.ParamMap["base_nickname"]; ok {
			for index, _ := range base_nicknames {
				recipe.BaseNicknames = append(recipe.BaseNicknames,
					semantic.NewString(section, cfg.Key("base_nickname"), semantic.WithLowercaseS(), semantic.OptsS(semantic.Index(index)), semantic.WithoutSpacesS()))
			}
		}

		if len(recipe.Product) == 0 {
			continue
		}

		frelconfig.Products = append(frelconfig.Products, recipe)

		for _, product := range recipe.Product {
			frelconfig.ProductsByNick[product.Nickname.Get()] = append(frelconfig.ProductsByNick[product.Nickname.Get()], recipe)

		}
	}

	return frelconfig

}

func (frelconfig *Config) Write() *file.File {
	inifile := frelconfig.Render()
	inifile.Write(inifile.File)
	return inifile.File
}
