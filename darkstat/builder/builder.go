package builder

import (
	"github.com/darklab8/fl-darkstat/darkstat/common/static_common"
	"github.com/darklab8/fl-darkstat/darkstat/front/static_front"
	"github.com/darklab8/go-utils/utils/timeit"
	"github.com/darklab8/go-utils/utils/utils_filepath"
	"github.com/darklab8/go-utils/utils/utils_types"
)

type Builder struct {
	components []*Component
	params     Params
}

type BuilderOption func(b *Builder)

type Params interface {
	GetBuildPath() utils_types.FilePath
}

func NewBuilder(params Params, opts ...BuilderOption) *Builder {
	b := &Builder{
		params: params,
	}
	for _, opt := range opts {
		opt(b)
	}
	return b
}

func (b *Builder) RegComps(components ...*Component) {
	b.components = append(b.components, components...)
}

func (b *Builder) build(components []*Component, params Params, filesystem *Filesystem) {

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
		target_folder := params.GetBuildPath().Join("static")
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

	b.build(b.components, b.params, filesystem)

	return filesystem
}
