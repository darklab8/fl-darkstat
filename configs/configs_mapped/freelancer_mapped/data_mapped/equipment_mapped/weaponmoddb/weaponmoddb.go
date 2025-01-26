package weaponmoddb

import (
	"github.com/darklab8/fl-darkstat/configs/cfg"
	"github.com/darklab8/fl-darkstat/configs/configs_mapped/parserutils/filefind/file"
	"github.com/darklab8/fl-darkstat/configs/configs_mapped/parserutils/iniload"
	"github.com/darklab8/fl-darkstat/configs/configs_mapped/parserutils/semantic"
)

const (
	FILENAME = "weaponmoddb.ini"
)

type ShieldMod struct {
	semantic.Model
	ShieldType     *semantic.String
	DamageModifier *semantic.Float
}

type WeaponType struct {
	semantic.Model
	Nickname   *semantic.String
	ShieldMods []*ShieldMod
}

type Config struct {
	*iniload.IniLoader

	WeaponTypes    []*WeaponType
	WeaponTypesMap map[string]*WeaponType
}

func Read(input_file *iniload.IniLoader) *Config {
	frelconfig := &Config{
		IniLoader:      input_file,
		WeaponTypesMap: make(map[string]*WeaponType),
	}
	if sections, ok := frelconfig.SectionMap["[weapontype]"]; ok {
		for _, section := range sections {
			weapon_type := &WeaponType{}
			weapon_type.Map(section)
			weapon_type.Nickname = semantic.NewString(section, cfg.Key("nickname"), semantic.WithLowercaseS(), semantic.WithoutSpacesS())

			KEY_SHIELD_MODE := cfg.Key("shield_mod")
			for index, _ := range section.ParamMap[KEY_SHIELD_MODE] {
				shield_mode := &ShieldMod{}
				shield_mode.Map(section)
				shield_mode.ShieldType = semantic.NewString(section, KEY_SHIELD_MODE, semantic.OptsS(semantic.Index(index)), semantic.WithLowercaseS(), semantic.WithoutSpacesS())
				shield_mode.DamageModifier = semantic.NewFloat(section, KEY_SHIELD_MODE, semantic.Precision(2), semantic.OptsF(semantic.Index(index), semantic.Order(1)))
				weapon_type.ShieldMods = append(weapon_type.ShieldMods, shield_mode)
			}
			frelconfig.WeaponTypes = append(frelconfig.WeaponTypes, weapon_type)
			frelconfig.WeaponTypesMap[weapon_type.Nickname.Get()] = weapon_type
		}
	}

	return frelconfig
}

func (frelconfig *Config) Write() *file.File {
	inifile := frelconfig.Render()
	inifile.Write(inifile.File)
	return inifile.File
}
