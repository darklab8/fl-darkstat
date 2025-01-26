package core_front

import (
	"embed"

	"github.com/darklab8/fl-darkstat/darkcore/core_types"
	"github.com/darklab8/fl-darkstat/darkcore/settings/logus"
	"github.com/darklab8/go-typelog/typelog"
	"github.com/darklab8/go-utils/utils/utils_types"
)

type StaticFilesystem struct {
	Files         []core_types.StaticFile
	relPathToFile map[utils_types.FilePath]core_types.StaticFile
}

func (fs StaticFilesystem) GetFileByRelPath(rel_path utils_types.FilePath) core_types.StaticFile {
	file, ok := fs.relPathToFile[rel_path]

	if !ok {
		logus.Log.Panic("expected file found by relpath", typelog.Any("relpath", rel_path))
	}

	return file
}

func GetFiles(fs embed.FS, params utils_types.GetFilesParams) StaticFilesystem {
	files := utils_types.GetFiles(fs, params)
	var filesystem StaticFilesystem = StaticFilesystem{
		relPathToFile: make(map[utils_types.FilePath]core_types.StaticFile),
	}

	for _, file := range files {
		var static_file_kind core_types.StaticFileKind

		switch file.Extension {
		case "js":
			static_file_kind = core_types.StaticFileJS
		case "css":
			static_file_kind = core_types.StaticFileCSS
		case "ico":
			static_file_kind = core_types.StaticFileIco
		}

		new_file := core_types.StaticFile{
			Filename: string(file.Relpath),
			Kind:     static_file_kind,
			Content:  string(file.Content),
		}
		filesystem.Files = append(filesystem.Files, new_file)
		filesystem.relPathToFile[file.Relpath] = new_file
	}
	return filesystem
}
