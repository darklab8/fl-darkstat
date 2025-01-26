package iniload

import (
	"github.com/darklab8/fl-darkstat/configs/configs_mapped/parserutils/filefind/file"
	"github.com/darklab8/fl-darkstat/configs/configs_mapped/parserutils/inireader"
	"github.com/darklab8/fl-darkstat/configs/configs_mapped/parserutils/semantic"
	"github.com/darklab8/fl-darkstat/configs/configs_settings/logus"
)

type IniLoader struct {
	semantic.ConfigModel
	input_file *file.File
	*inireader.INIFile
}

func NewLoader(input_file *file.File) *IniLoader {
	fileconfig := &IniLoader{input_file: input_file}
	return fileconfig
}

// Scan is heavy operations for goroutine ^_^
func (fileconfig *IniLoader) Scan() *IniLoader {
	if fileconfig.input_file == nil {
		logus.Log.Error("input_file is empty")
	}

	iniconfig := inireader.Read(fileconfig.input_file)
	fileconfig.Init(iniconfig.Sections, iniconfig.Comments, iniconfig.File.GetFilepath())
	fileconfig.INIFile = iniconfig
	return fileconfig
}
