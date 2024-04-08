package initialworld

import (
	"github.com/darklab8/fl-configs/configs/configs_mapped/parserutils/filefind/file"
	"github.com/darklab8/fl-configs/configs/configs_mapped/parserutils/iniload"
	"github.com/darklab8/fl-configs/configs/configs_mapped/parserutils/semantic"
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
}

func Read(input_file *iniload.IniLoader) *Config {
	frelconfig := &Config{
		IniLoader: input_file,
		Groups:    make([]*Group, 0, 100),
		GroupsMap: make(map[string]*Group),
	}

	if groups, ok := frelconfig.SectionMap["[Group]"]; ok {

		for _, group_res := range groups {
			group := &Group{}
			group.Map(group_res)
			group.Nickname = semantic.NewString(group_res, "nickname", semantic.WithLowercaseS(), semantic.WithoutSpacesS())
			group.IdsName = semantic.NewInt(group_res, "ids_name")
			group.IdsInfo = semantic.NewInt(group_res, "ids_info")
			group.IdsShortName = semantic.NewInt(group_res, "ids_short_name")

			group.Relationships = make([]*Relationship, 0, 20)

			param_rep_key := "rep"
			for rep_index, _ := range group_res.ParamMap[param_rep_key] {

				rep := &Relationship{}
				rep.Map(group_res)
				rep.Rep = semantic.NewFloat(group_res, param_rep_key, semantic.Precision(2), semantic.Index(rep_index))
				rep.TargetNickname = semantic.NewString(group_res, param_rep_key, semantic.OptsS(semantic.Index(rep_index), semantic.Order(1)), semantic.WithLowercaseS(), semantic.WithoutSpacesS())
				group.Relationships = append(group.Relationships, rep)
			}

			frelconfig.Groups = append(frelconfig.Groups, group)
			frelconfig.GroupsMap[group.Nickname.Get()] = group
		}
	}

	return frelconfig
}

func (frelconfig *Config) Write() *file.File {
	inifile := frelconfig.Render()
	inifile.Write(inifile.File)
	return inifile.File
}
