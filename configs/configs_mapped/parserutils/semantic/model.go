package semantic

import (
	"github.com/darklab8/fl-darkstat/configs/configs_mapped/parserutils/filefind/file"
	"github.com/darklab8/fl-darkstat/configs/configs_mapped/parserutils/inireader"

	"github.com/darklab8/go-utils/utils/utils_types"
)

type Model struct {
	section *inireader.Section
}

func (s *Model) Map(section *inireader.Section) {
	s.section = section
}

func (s *Model) RenderModel() *inireader.Section {
	return s.section
}

type ConfigModel struct {
	sections []*inireader.Section
	comments []string
	filepath utils_types.FilePath
}

func (s *ConfigModel) Init(sections []*inireader.Section, comments []string, filepath utils_types.FilePath) {
	s.sections = sections
	s.comments = comments
	s.filepath = filepath
}

func (s *ConfigModel) SetOutputPath(filepath utils_types.FilePath) {
	s.filepath = filepath
}

func (s *ConfigModel) Render() *inireader.INIFile {
	inifile := &inireader.INIFile{}
	inifile.Comments = s.comments
	inifile.Sections = s.sections
	inifile.File = file.NewFile(s.filepath)
	return inifile
}
