package empathy_mapped

import (
	"github.com/darklab8/fl-darkstat/configs/cfg"
	"github.com/darklab8/fl-darkstat/configs/configs_mapped/parserutils/filefind/file"
	"github.com/darklab8/fl-darkstat/configs/configs_mapped/parserutils/iniload"
	"github.com/darklab8/fl-darkstat/configs/configs_mapped/parserutils/semantic"

	"github.com/darklab8/go-utils/utils/utils_types"
)

type EmpathyRate struct {
	semantic.Model

	TargetFactionNickname *semantic.String // 0
	RepoChange            *semantic.Float  // 1
}

type RepChangeEffects struct {
	semantic.Model
	Group *semantic.String

	ObjectDestruction *semantic.Float
	MissionSuccess    *semantic.Float
	MissionFailure    *semantic.Float
	MissionAbort      *semantic.Float

	EmpathyRates    []*EmpathyRate
	EmpathyRatesMap map[string]*EmpathyRate
}

type Config struct {
	*iniload.IniLoader
	RepChangeEffects []*RepChangeEffects
	RepoChangeMap    map[string]*RepChangeEffects
}

const (
	FILENAME utils_types.FilePath = "empathy.ini"
)

func Read(input_file *iniload.IniLoader) *Config {
	frelconfig := &Config{IniLoader: input_file}

	frelconfig.RepChangeEffects = make([]*RepChangeEffects, 0, 20)
	frelconfig.RepoChangeMap = make(map[string]*RepChangeEffects)

	for _, section := range input_file.SectionMap["[repchangeeffects]"] {
		repo_changes := &RepChangeEffects{}
		repo_changes.Map(section)
		repo_changes.Group = semantic.NewString(section, cfg.Key("group"), semantic.WithLowercaseS(), semantic.WithoutSpacesS())
		repo_changes.EmpathyRatesMap = make(map[string]*EmpathyRate)

		event_key := cfg.Key("event")
		for event_index, event := range section.ParamMap[event_key] {
			switch event.First.AsString() {
			case "object_destruction":
				repo_changes.ObjectDestruction = semantic.NewFloat(section, event_key, semantic.Precision(2), semantic.OptsF(semantic.Index(event_index), semantic.Order(1)))
			case "random_mission_success":
				repo_changes.MissionSuccess = semantic.NewFloat(section, event_key, semantic.Precision(2), semantic.OptsF(semantic.Index(event_index), semantic.Order(1)))
			case "random_mission_failure":
				repo_changes.MissionFailure = semantic.NewFloat(section, event_key, semantic.Precision(2), semantic.OptsF(semantic.Index(event_index), semantic.Order(1)))
			case "random_mission_abortion":
				repo_changes.MissionAbort = semantic.NewFloat(section, event_key, semantic.Precision(2), semantic.OptsF(semantic.Index(event_index), semantic.Order(1)))
			}
		}

		empathy_rate_key := cfg.Key("empathy_rate")
		for good_index, _ := range section.ParamMap[empathy_rate_key] {
			empathy := &EmpathyRate{}
			empathy.Map(section)
			empathy.TargetFactionNickname = semantic.NewString(section, empathy_rate_key, semantic.OptsS(semantic.Index(good_index), semantic.Order(0)))
			empathy.RepoChange = semantic.NewFloat(section, empathy_rate_key, semantic.Precision(2), semantic.OptsF(semantic.Index(good_index), semantic.Order(1)))
			repo_changes.EmpathyRates = append(repo_changes.EmpathyRates, empathy)
			repo_changes.EmpathyRatesMap[empathy.TargetFactionNickname.Get()] = empathy
		}

		frelconfig.RepChangeEffects = append(frelconfig.RepChangeEffects, repo_changes)
		frelconfig.RepoChangeMap[repo_changes.Group.Get()] = repo_changes
	}
	return frelconfig

}

func (frelconfig *Config) Write() *file.File {
	inifile := frelconfig.Render()
	inifile.Write(inifile.File)
	return inifile.File
}
