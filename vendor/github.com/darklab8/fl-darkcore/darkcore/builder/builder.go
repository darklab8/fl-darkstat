package builder

import (
	"github.com/darklab8/go-utils/utils/timeit"
	"github.com/darklab8/go-utils/utils/utils_filepath"
	"github.com/darklab8/go-utils/utils/utils_types"
)

type Builder struct {
	components   []*Component
	params       Params
	static_files []StaticFile
}

type StaticFile struct {
	path    utils_types.FilePath
	content []byte
}

func NewStaticFile(path utils_types.FilePath, content []byte) StaticFile {
	return StaticFile{
		path:    path,
		content: content,
	}
}

type BuilderOption func(b *Builder)

type Params interface {
	GetBuildPath() utils_types.FilePath
}

func NewBuilder(params Params, static_files []StaticFile, opts ...BuilderOption) *Builder {
	b := &Builder{
		params:       params,
		static_files: static_files,
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

		for _, static_file := range b.static_files {
			filesystem.WriteToMem(utils_filepath.Join(target_folder, static_file.path), []byte(static_file.content))
		}

	}, timeit.WithMsg("gathered static assets"))
}

func (b *Builder) BuildAll() *Filesystem {
	build_root := utils_types.FilePath("build")
	filesystem := NewFileystem(build_root)

	b.build(b.components, b.params, filesystem)

	return filesystem
}
