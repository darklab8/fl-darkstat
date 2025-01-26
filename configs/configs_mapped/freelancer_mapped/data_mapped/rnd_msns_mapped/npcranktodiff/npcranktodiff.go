package npcranktodiff

import (
	"github.com/darklab8/fl-darkstat/configs/cfg"
	"github.com/darklab8/fl-darkstat/configs/configs_mapped/parserutils/filefind/file"
	"github.com/darklab8/fl-darkstat/configs/configs_mapped/parserutils/iniload"
	"github.com/darklab8/fl-darkstat/configs/configs_mapped/parserutils/semantic"
	"github.com/darklab8/go-utils/utils/utils_types"
)

type NPCRankToDifficulty struct {
	semantic.Model
	Rank         *semantic.Int
	Difficulties []*semantic.Float
}

type Config struct {
	*iniload.IniLoader
	NPCRankToDifficulties []*NPCRankToDifficulty
}

const (
	FILENAME utils_types.FilePath = "npcranktodiff.ini"
)

func Read(input_file *iniload.IniLoader) *Config {
	frelconfig := &Config{IniLoader: input_file}

	for _, section := range input_file.SectionMap["[rankandformationsizetodifficulty]"] {

		for index, values := range section.ParamMap[cfg.Key("npcrank")] {
			npc_rank_to_diff := &NPCRankToDifficulty{
				Rank: semantic.NewInt(section, cfg.Key("npcrank"), semantic.Index(index), semantic.Order(0)),
			}

			len_of_difficulties := len(values.Values) - 1
			for i := 1; i <= len_of_difficulties; i++ {
				npc_rank_to_diff.Difficulties = append(npc_rank_to_diff.Difficulties,
					semantic.NewFloat(section, cfg.Key("npcrank"), semantic.Precision(2), semantic.OptsF(semantic.Index(index), semantic.Order(i))),
				)
			}

			npc_rank_to_diff.Map(section)
			frelconfig.NPCRankToDifficulties = append(frelconfig.NPCRankToDifficulties, npc_rank_to_diff)
		}

	}
	return frelconfig

}

func (frelconfig *Config) Write() *file.File {
	inifile := frelconfig.Render()
	inifile.Write(inifile.File)
	return inifile.File
}
