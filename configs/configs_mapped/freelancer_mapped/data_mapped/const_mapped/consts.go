package const_mapped

import (
	"github.com/darklab8/fl-darkstat/configs/cfg"
	"github.com/darklab8/fl-darkstat/configs/configs_mapped/parserutils/filefind/file"
	"github.com/darklab8/fl-darkstat/configs/configs_mapped/parserutils/iniload"
	"github.com/darklab8/fl-darkstat/configs/configs_mapped/parserutils/semantic"
)

const (
	FILENAME = "constants.ini"
)

type ShieldEquipConsts struct {
	semantic.Model
	HULL_DAMAGE_FACTOR *semantic.Float
}

type EngineEquipConsts struct {
	semantic.Model
	CRUISING_SPEED *semantic.Int
}

type Config struct {
	*iniload.IniLoader

	ShieldEquipConsts *ShieldEquipConsts
	EngineEquipConsts *EngineEquipConsts
}

func Read(input_file *iniload.IniLoader) *Config {
	frelconfig := &Config{
		IniLoader: input_file,
	}
	if groups, ok := frelconfig.SectionMap["[shieldequipconsts]"]; ok {
		shield_consts := &ShieldEquipConsts{}
		shield_consts.Map(groups[0])
		shield_consts.HULL_DAMAGE_FACTOR = semantic.NewFloat(groups[0], cfg.Key("hull_damage_factor"), semantic.Precision(2))

		frelconfig.ShieldEquipConsts = shield_consts
	}
	if groups, ok := frelconfig.SectionMap["[engineequipconsts]"]; ok {
		const_group := &EngineEquipConsts{}
		const_group.Map(groups[0])
		const_group.CRUISING_SPEED = semantic.NewInt(groups[0], cfg.Key("cruising_speed"))

		frelconfig.EngineEquipConsts = const_group
	}

	return frelconfig
}

func (frelconfig *Config) Write() *file.File {
	inifile := frelconfig.Render()
	inifile.Write(inifile.File)
	return inifile.File
}
