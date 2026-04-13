package builder

import (
	"context"
	"fmt"
	"time"

	"github.com/darklab8/fl-darkstat/darkcore/core_types"
	"github.com/darklab8/fl-darkstat/darkcore/settings"
	darkstat_settings "github.com/darklab8/fl-darkstat/darkstat/settings"

	"github.com/darklab8/go-utils/utils/timeit"
	"github.com/darklab8/go-utils/utils/utils_filepath"
	"github.com/darklab8/go-utils/utils/utils_types"
)

type Builder struct {
	components   []*Component
	params       Params
	static_files []StaticFile
}

func (b *Builder) GetStaticFileChecker() map[utils_types.FilePath]bool {
	checker := make(map[utils_types.FilePath]bool)
	for _, file := range b.static_files {
		checker[file.path] = true
	}
	return checker
}

func (b *Builder) GetParams() Params { return b.params }

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

func NewStaticFileFromCore(s core_types.StaticFile) StaticFile {
	return NewStaticFile(utils_types.FilePath(s.Filename), []byte(s.Content))
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

func (b *Builder) AddStaticFiles(static_files []StaticFile) {
	b.static_files = append(b.static_files, static_files...)
}

func (b *Builder) RegComps(components ...*Component) {
	b.components = append(b.components, components...)
}

func chunkSlice(slice []*Component, chunkSize int) [][]*Component {
	var chunks [][]*Component
	for i := 0; i < len(slice); i += chunkSize {
		end := i + chunkSize

		// necessary check to avoid slicing beyond
		// slice capacity
		if end > len(slice) {
			end = len(slice)
		}

		chunks = append(chunks, slice[i:end])
	}

	return chunks
}

// func (b *Builder) ToWebServer() *Filesystem {
// }

type CleanFolderKind int8

const (
	NotCleanFolder CleanFolderKind = iota
	YesCleanFolder
)

func (b *Builder) BuildAll(to_mem bool, cleanup_build_folder CleanFolderKind, filesystem *Filesystem) *Filesystem {
	var ctx context.Context = context.Background()

	if !to_mem {
		darkstat_settings.Env.IsStaticSiteGenerator = true
	}

	build_root := utils_types.FilePath(settings.Env.BuildFolder)
	if filesystem == nil {
		filesystem = NewFileystem(build_root)
	}

	filesystem.CreateBuildFolder(cleanup_build_folder)
	fmt.Println("beginning build operation")
	results := make(chan WriteResult)

	time_start := time.Now()

	timeit.NewTimerF(func() {
		chunked_components := chunkSlice(b.components, 10000)
		len_comps := len(chunked_components)
		fmt.Println("components chunks", len_comps)
		for chunk_index, components_chunk := range chunked_components {

			if to_mem {
				for _, comp := range components_chunk {
					filesystem.WriteToMem(comp.GetPagePath(b.params), &MemComp{
						comp: comp,
						b:    b,
					})
				}
			} else {
				type ModeToRun int8
				const (
					ModeNormal     ModeToRun = iota
					ModeGoroutines           // same time as mode normal
					ModeWorkerPool           // same time as mode normal.
					// i guess filesystem is bottleneck
					// leaving code for other modes here just for inspiration and history

					// i'll leave mode normal then
				)

				mode_to_run := ModeNormal
				switch mode_to_run {
				case ModeWorkerPool:
					worker := func(id int, jobs <-chan *Component, results chan<- WriteResult) {
						// fmt.Println("worker", id, "started  processing jobs")
						for comp := range jobs {
							results <- comp.Write(ctx, b.params)
						}
						// fmt.Println("worker", id, "finished  processing jobs")
					}
					numJobs := len(components_chunk)
					jobs := make(chan *Component, numJobs)
					results := make(chan WriteResult, numJobs)

					for w := 0; w <= 3; w++ {
						go worker(w, jobs, results)
					}
					for _, comp := range components_chunk {
						jobs <- comp
					}
					close(jobs)

					for range components_chunk {
						result := <-results
						filesystem.WriteToFile(result.realpath, result.bytes)
					}
				case ModeGoroutines:
					for _, comp := range components_chunk {
						go func(comp *Component) {
							results <- comp.Write(ctx, b.params)
						}(comp)
					}
					for range components_chunk {
						result := <-results
						filesystem.WriteToFile(result.realpath, result.bytes)
					}
				default: // normal

					for _, comp := range components_chunk {
						result := comp.Write(ctx, b.params)
						filesystem.WriteToFile(result.realpath, result.bytes)
					}
				}

			}

			if chunk_index%10 == 0 {
				fmt.Println("finished chunk=", chunk_index, "/", len_comps, " time=", time.Since(time_start))
			}
		}

	}, timeit.WithMsg("wrote components"))

	timeit.NewTimerF(func() {
		target_folder := b.params.GetBuildPath().Join("static")
		for _, static_file := range b.static_files {
			path := utils_filepath.Join(target_folder, static_file.path)
			if to_mem {
				filesystem.WriteToMem(path, &MemStatic{
					content: static_file.content,
				})
			} else {
				filesystem.WriteToFile(path, []byte(static_file.content))
			}
		}
	}, timeit.WithMsg("gathered static assets"))

	return filesystem
}
