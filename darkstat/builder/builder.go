package builder

import (
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"github.com/darklab8/fl-darkstat/darkstat/common/types"
	"github.com/darklab8/fl-darkstat/darkstat/settings"
	"github.com/darklab8/fl-darkstat/darkstat/settings/logus"
	"github.com/darklab8/go-utils/goutils/utils/utils_filepath"
	"github.com/darklab8/go-utils/goutils/utils/utils_os"
	"github.com/darklab8/go-utils/goutils/utils/utils_types"
)

type Builder struct {
	components []*Component
}

type BuilderOption func(b *Builder)

func NewBuilder(opts ...BuilderOption) *Builder {
	b := &Builder{}
	for _, opt := range opts {
		opt(b)
	}
	return b
}

func (b *Builder) RegComps(components ...*Component) {
	b.components = append(b.components, components...)
}

func (b *Builder) build(params types.GlobalParams, filesystem *Filesystem) {

	for _, comp := range b.components {
		comp.Write(params, filesystem)
	}

	folders := utils_os.GetRecursiveDirs(settings.ProjectFolder)
	for _, folder := range folders {
		if utils_filepath.Base(folder) == "static" {

			filepath.WalkDir(folder.ToString(), func(path string, d fs.DirEntry, err error) error {
				if logus.Log.CheckError(err, "err is encountered during filepath.Walkdir") {
					return nil
				}
				if d.IsDir() {
					return nil
				}

				target_folder := utils_filepath.Join(utils_types.FilePath(params.Buildpath.ToString()), "static").ToString()

				data, err := os.ReadFile(path)
				if logus.Log.CheckError(err, "failed to read file at path in WalkDir") {
					return nil
				}

				new_path := strings.Replace(path, folder.ToString(), target_folder, 1)

				filesystem.WriteToMem(utils_types.FilePath(new_path), data)
				return nil
			})
		}
	}
}

func (b *Builder) BuildAll() *Filesystem {
	build_root := utils_types.FilePath("build")
	filesystem := NewFileystem(build_root)

	staticPrefix := "static/"

	var siteRoot string
	if value, ok := os.LookupEnv("SITE_ROOT"); ok {
		siteRoot = value
	} else {
		siteRoot = "/"
	}
	b.build(types.GlobalParams{
		Buildpath:         "",
		Theme:             types.ThemeLight,
		SiteRoot:          siteRoot,
		StaticRoot:        siteRoot + staticPrefix,
		OppositeThemeRoot: siteRoot + "dark/",
	}, filesystem)
	b.build(types.GlobalParams{
		Buildpath:         utils_filepath.Join("dark"),
		Theme:             types.ThemeDark,
		SiteRoot:          siteRoot + "dark/",
		StaticRoot:        siteRoot + "dark/" + staticPrefix,
		OppositeThemeRoot: siteRoot,
	}, filesystem)

	return filesystem
}
