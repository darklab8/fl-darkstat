package builder

import (
	"time"

	"github.com/darklab8/fl-darkstat/darkstat/common/static_common"
	"github.com/darklab8/fl-darkstat/darkstat/common/types"
	"github.com/darklab8/fl-darkstat/darkstat/front/static_front"
	"github.com/darklab8/fl-darkstat/darkstat/settings"
	"github.com/darklab8/go-utils/utils/timeit"
	"github.com/darklab8/go-utils/utils/utils_filepath"
	"github.com/darklab8/go-utils/utils/utils_types"
)

type Builder struct {
	components []*Component
	dark_pages []*Component
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

func (b *Builder) RegDark(components ...*Component) {
	b.dark_pages = append(b.dark_pages, components...)
}

func (b *Builder) build(components []*Component, params types.GlobalParams, filesystem *Filesystem) {

	timeit.NewTimerF(func(m *timeit.Timer) {
		results := make(chan WriteResult)
		for _, comp := range components {
			go func(comp *Component) {
				results <- comp.Write(params)
			}(comp)
		}
		for range components {
			result := <-results
			filesystem.WriteToMem(result.realpath, result.bytes)
		}
	}, timeit.WithMsg("wrote components"))

	timeit.NewTimerF(func(m *timeit.Timer) {
		target_folder := utils_filepath.Join(utils_types.FilePath(params.Buildpath.ToString()), "static")
		filesystem.WriteToMem(utils_filepath.Join(target_folder, "htmx.js"), []byte(static_front.HtmxJs))
		filesystem.WriteToMem(utils_filepath.Join(target_folder, "preload.js"), []byte(static_front.PreloadJs))
		filesystem.WriteToMem(utils_filepath.Join(target_folder, "sortable.js"), []byte(static_front.SortableJs))
		filesystem.WriteToMem(utils_filepath.Join(target_folder, "custom.js"), []byte(static_front.CustomJS))
		filesystem.WriteToMem(utils_filepath.Join(target_folder, "common", "favicon.ico"), []byte(static_common.FaviconIco))
	}, timeit.WithMsg("gathered static assets"))
}

func (b *Builder) BuildAll() *Filesystem {
	build_root := utils_types.FilePath("build")
	filesystem := NewFileystem(build_root)

	staticPrefix := "static/"

	siteRoot := settings.Env.SiteRoot
	b.build(b.components, types.GlobalParams{
		Buildpath:         "",
		Theme:             types.ThemeLight,
		SiteRoot:          siteRoot,
		StaticRoot:        siteRoot + staticPrefix,
		OppositeThemeRoot: siteRoot + "dark.html",
		Heading:           settings.Env.AppHeading,
		Timestamp:         time.Now().UTC(),
	}, filesystem)

	// // Implement dark theme later
	// // u need only Index page rebuilded, not all of them ^_^
	// b.build(b.dark_pages, types.GlobalParams{
	// 	Buildpath:         "",
	// 	Theme:             types.ThemeDark,
	// 	SiteRoot:          siteRoot,
	// 	StaticRoot:        siteRoot + staticPrefix,
	// 	OppositeThemeRoot: siteRoot,
	// 	Heading:           settings.Conf.AppHeading,
	// }, filesystem)

	return filesystem
}
