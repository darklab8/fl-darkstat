package diff2money

import (
	"github.com/darklab8/fl-configs/configs/configs_mapped/parserutils/filefind/file"
	"github.com/darklab8/fl-configs/configs/configs_mapped/parserutils/iniload"
	"github.com/darklab8/fl-configs/configs/configs_mapped/parserutils/semantic"
	"github.com/darklab8/go-utils/goutils/utils/utils_types"
)

type DiffToMoney struct {
	semantic.Model
	MinLevel   *semantic.Float
	MoneyAward *semantic.Int
}

type Config struct {
	*iniload.IniLoader
	DiffToMoney []*DiffToMoney
}

const (
	FILENAME utils_types.FilePath = "diff2money.ini"
)

func Read(input_file *iniload.IniLoader) *Config {
	frelconfig := &Config{IniLoader: input_file}

	for _, section := range input_file.SectionMap["[Diff2Money]"] {

		for index, _ := range section.ParamMap["diff2money"] {
			diff_to_money := &DiffToMoney{
				MinLevel:   semantic.NewFloat(section, "diff2money", semantic.Precision(2), semantic.Index(index), semantic.Order(0)),
				MoneyAward: semantic.NewInt(section, "diff2money", semantic.Index(index), semantic.Order(1)),
			}
			diff_to_money.Map(section)
			frelconfig.DiffToMoney = append(frelconfig.DiffToMoney, diff_to_money)
		}

	}
	return frelconfig

}

func (frelconfig *Config) Write() *file.File {
	inifile := frelconfig.Render()
	inifile.Write(inifile.File)
	return inifile.File
}
