package initialworld

import (
	"strconv"

	"github.com/darklab8/fl-darkstat/configs/cfg"
	"github.com/darklab8/fl-darkstat/configs/configs_mapped/freelancer_mapped/data_mapped/initialworld/flhash"
	"github.com/darklab8/fl-darkstat/configs/configs_mapped/parserutils/filefind/file"
	"github.com/darklab8/fl-darkstat/configs/configs_mapped/parserutils/iniload"
	"github.com/darklab8/fl-darkstat/configs/configs_mapped/parserutils/semantic"
	"github.com/darklab8/fl-darkstat/configs/configs_settings/logus"
)

const (
	FILENAME = "initialworld.ini"
)

type Relationship struct {
	semantic.Model

	Rep            *semantic.Float
	TargetNickname *semantic.String
}

type Group struct {
	semantic.Model

	Nickname *semantic.String
	IdsName  *semantic.Int
	IdsInfo  *semantic.Int

	IdsShortName  *semantic.Int
	Relationships []*Relationship
}

type Config struct {
	*iniload.IniLoader

	Groups    []*Group
	GroupsMap map[string]*Group

	LockedGates map[flhash.HashCode]bool
}

func Read(input_file *iniload.IniLoader) *Config {
	config := &Config{
		IniLoader:   input_file,
		Groups:      make([]*Group, 0, 100),
		GroupsMap:   make(map[string]*Group),
		LockedGates: make(map[flhash.HashCode]bool),
	}

	if locked_gates, ok := config.SectionMap["[locked_gates]"]; ok {

		for _, values := range locked_gates[0].ParamMap[cfg.Key("locked_gate")] {

			int_hash, err := strconv.Atoi(values.First.AsString())

			logus.Log.CheckPanic(err, "failed to parse locked_gate")

			config.LockedGates[flhash.HashCode(int_hash)] = true
		}
	}

	if groups, ok := config.SectionMap["[group]"]; ok {

		for _, group_res := range groups {
			group := &Group{}
			group.Map(group_res)
			group.Nickname = semantic.NewString(group_res, cfg.Key("nickname"), semantic.WithLowercaseS(), semantic.WithoutSpacesS())
			group.IdsName = semantic.NewInt(group_res, cfg.Key("ids_name"), semantic.Optional())
			group.IdsInfo = semantic.NewInt(group_res, cfg.Key("ids_info"), semantic.Optional())
			group.IdsShortName = semantic.NewInt(group_res, cfg.Key("ids_short_name"))

			group.Relationships = make([]*Relationship, 0, 20)

			param_rep_key := cfg.Key("rep")
			for rep_index, _ := range group_res.ParamMap[param_rep_key] {

				rep := &Relationship{}
				rep.Map(group_res)
				rep.Rep = semantic.NewFloat(group_res, param_rep_key, semantic.Precision(2), semantic.OptsF(semantic.Index(rep_index)))
				rep.TargetNickname = semantic.NewString(group_res, param_rep_key, semantic.OptsS(semantic.Index(rep_index), semantic.Order(1)), semantic.WithLowercaseS(), semantic.WithoutSpacesS())
				group.Relationships = append(group.Relationships, rep)
			}

			config.Groups = append(config.Groups, group)
			config.GroupsMap[group.Nickname.Get()] = group
		}
	}

	return config
}

func (frelconfig *Config) Write() *file.File {
	inifile := frelconfig.Render()
	inifile.Write(inifile.File)
	return inifile.File
}
