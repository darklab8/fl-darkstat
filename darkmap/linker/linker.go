package linker

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/darklab8/fl-darkstat/darkcore/builder"
	"github.com/darklab8/fl-darkstat/darkcore/core_static"
	"github.com/darklab8/fl-darkstat/darkcore/core_types"

	"github.com/darklab8/fl-darkstat/darkmap/export_map"
	"github.com/darklab8/fl-darkstat/darkmap/front"
	"github.com/darklab8/fl-darkstat/darkmap/front/static"
	"github.com/darklab8/fl-darkstat/darkmap/front/static_front"
	"github.com/darklab8/fl-darkstat/darkmap/front/urls"
	"github.com/darklab8/fl-darkstat/darkmap/settings"
	"github.com/darklab8/fl-darkstat/darkmap/settings/logus"
	"github.com/darklab8/fl-darkstat/darkmap/types"
	"github.com/darklab8/fl-darkstat/darkmap/utfextract"
	"github.com/darklab8/go-utils/typelog"
	"github.com/darklab8/go-utils/utils/timeit"
	"github.com/darklab8/go-utils/utils/utils_types"
)

type Linker struct {
	Export   *export_map.Export
	ToMemory bool
}

type LinkOption func(l *Linker)

func NewLinker(to_memory bool, opts ...LinkOption) *Linker {
	l := &Linker{
		ToMemory: to_memory,
	}
	for _, opt := range opts {
		opt(l)
	}

	return l
}

func (l *Linker) Link(ctx context.Context) *builder.Builder {
	l.Export = export_map.NewExport(ctx)

	defer timeit.NewTimer("Link").Close()
	var build *builder.Builder
	staticPrefix := "static/"
	siteRoot := settings.Env.SiteRoot
	params := types.GlobalParams{
		Buildpath:         "",
		SiteRoot:          siteRoot,
		StaticRoot:        siteRoot + staticPrefix,
		OppositeThemeRoot: siteRoot + "dark.html",
		Timestamp:         time.Now().UTC(),
	}

	files := []builder.StaticFile{
		builder.NewStaticFileFromCore(core_static.HtmxJS),
		builder.NewStaticFileFromCore(core_static.HtmxPreloadJS),
		builder.NewStaticFileFromCore(core_static.SortableJS),
		builder.NewStaticFileFromCore(core_static.ResetCSS),
		builder.NewStaticFileFromCore(static_front.FaviconIco),
		builder.NewStaticFileFromCore(static_front.CommonCSS),
		builder.NewStaticFileFromCore(static_front.CustomCSS),
		builder.NewStaticFileFromCore(static_front.GalaxyCSS),
		builder.NewStaticFileFromCore(static_front.CustomJS),
		builder.NewStaticFileFromCore(static_front.MapGalaxyJS),
		builder.NewStaticFileFromCore(static_front.MapSystemJS),
		builder.NewStaticFileFromCore(static_front.PanzoomJS),
	}

	for _, file := range static.StaticFilesystem.Files {
		files = append(files, builder.NewStaticFileFromCore(file))
	}
	build = builder.NewBuilder(params, files)

	build.RegComps(
		builder.NewComponent(
			urls.Index,
			front.Index(l.Export),
		),
	)

	for _, system := range l.Export.Systems {
		build.RegComps(
			builder.NewComponent(
				utils_types.FilePath(front.SystemDetailedUrl(system)),
				front.System(system),
			),
		)
	}

	var extra_files []builder.StaticFile

	type StaticFileInParallel struct {
		core_types.StaticFile
		made bool
	}
	decoded_shape_files := make(chan StaticFileInParallel)

	time_start := time.Now()
	fmt.Println("SHAPES STARTING SENDING JOBS")
	created_jobs := 0
	for _, shape := range l.Export.Shapes.ShapesByNick {
		if l.ToMemory {
			if _, permitted := l.Export.Shapes.PermittedShapes[strings.ToLower(shape.Nickname)]; !permitted {
				continue
			}
		}

		created_jobs++
		go func() {

			var result StaticFileInParallel
			if shape.Nickname == "dsy_planet_earthgrncld" {
				fmt.Print()
			}

			image, err := SelectImage(shape, "tga")
			if err != nil {
				image, err = SelectImage(shape, "dds")
			}
			if logus.Log.CheckWarn(err, "not found image for shape, skipping",
				typelog.Any("shape", shape.Nickname+"."+shape.Extension),
			) {
				decoded_shape_files <- result
				return
			}
			jpeg_result, err := utfextract.TransformToJpeg(image)
			if logus.Log.CheckWarn(err, fmt.Sprintln("unable decoding "+image.Extension+" image"),
				typelog.Any("image_name", image.Nickname+"."+image.Extension),
				typelog.Any("shape", shape.Nickname+"."+shape.Extension),
				typelog.Any("image_dest", image.Dest),
			) {
				decoded_shape_files <- result
				return
			}
			result.StaticFile = core_types.StaticFile{
				Content:  jpeg_result.String(),
				Filename: fmt.Sprintf("%s.png", shape.Nickname),
				Kind:     core_types.StaticFileUnknown,
			}
			result.made = true
			decoded_shape_files <- result
		}()
	}
	fmt.Println("SHAPES SENT ALL JOBS")

	for i := 0; i < created_jobs; i++ {
		result := <-decoded_shape_files
		if !result.made {
			continue
		}
		extra_files = append(extra_files, builder.NewStaticFileFromCore(result.StaticFile))
	}
	fmt.Println("SHAPES FINISHED ACCEPTING ALL JOBS, took time seconds=", time.Since(time_start).Seconds(), " handled jobs=", created_jobs)

	build.AddStaticFiles(extra_files)

	return build
}

func SelectImage(shape *utfextract.Shape, extension string) (*utfextract.Image, error) {

	var result *utfextract.Image

	for _, image := range shape.Images {
		if image.Extension == extension {
			if result == nil {
				result = image
			} else {
				if len(image.Data) > len(result.Data) {
					result = image
				}
			}

		}
	}

	if result != nil {
		return result, nil
	}

	return nil, errors.New(fmt.Sprintln("not found image ", extension))
}
