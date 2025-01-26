package interface_mapped

import (
	"strconv"

	"github.com/darklab8/fl-darkstat/configs/cfg"
	"github.com/darklab8/fl-darkstat/configs/configs_mapped/parserutils/filefind/file"
	"github.com/darklab8/fl-darkstat/configs/configs_mapped/parserutils/iniload"
	"github.com/darklab8/fl-darkstat/configs/configs_mapped/parserutils/inireader/inireader_types"
	"github.com/darklab8/fl-darkstat/configs/configs_mapped/parserutils/semantic"
	"github.com/darklab8/fl-darkstat/configs/configs_settings/logus"
)

const (
	FILENAME_FL_INI                                     = "infocardmap.ini"
	RESOURCE_HEADER_MAP_TABLE inireader_types.IniHeader = "[infocardmaptable]"
	RESOURCE_KEY_MAP                                    = "map"
)

type InfocardMapTable struct {
	semantic.Model
	Map map[int]int
}

type Config struct {
	*iniload.IniLoader
	InfocardMapTable InfocardMapTable
}

func Read(input_file *iniload.IniLoader) *Config {
	frelconfig := &Config{
		IniLoader:        input_file,
		InfocardMapTable: InfocardMapTable{Map: make(map[int]int)},
	}

	if resources, ok := input_file.SectionMap[RESOURCE_HEADER_MAP_TABLE]; ok {

		for _, mappy := range resources[0].ParamMap[cfg.Key("map")] {

			id_key, err := strconv.Atoi(mappy.First.AsString())
			logus.Log.CheckPanic(err, "failed to read number from infocardmaptable")

			id_value, err := strconv.Atoi(mappy.Values[1].AsString())
			logus.Log.CheckPanic(err, "failed to read number from infocardmaptable")

			frelconfig.InfocardMapTable.Map[id_key] = id_value
		}
	}

	return frelconfig
}

func (frelconfig *Config) Write() *file.File {
	inifile := frelconfig.Render()
	inifile.Write(inifile.File)
	return inifile.File
}
