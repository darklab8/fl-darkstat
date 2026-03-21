package linker

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/darklab8/fl-darkstat/darkcore/builder"
	"github.com/darklab8/fl-darkstat/darkcore/core_static"
	"github.com/darklab8/fl-darkstat/darkcore/core_types"

	"github.com/darklab8/fl-darkstat/darkmap/front"
	"github.com/darklab8/fl-darkstat/darkmap/front/export_map"
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
	Export *export_map.Export
}

type LinkOption func(l *Linker)

func NewLinker(opts ...LinkOption) *Linker {
	l := &Linker{}
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

	for _, shape := range l.Export.Shapes.ShapesByNick {

		if shape.Nickname == "nav_terrain_ice" {
			fmt.Print()
		}

		image, err := SelectImageTga(shape)
		if logus.Log.CheckWarn(err, "not found imega tga for shape, skipping",
			typelog.Any("shape", shape.Nickname+"."+shape.Extension),
		) {
			continue
		}
		jpeg_result, err := utfextract.TransformToJpeg(image.Data)
		if logus.Log.CheckWarn(err, "unable decoding tga image",
			typelog.Any("image_name", image.Nickname+"."+image.Extension),
			typelog.Any("shape", shape.Nickname+"."+shape.Extension),
		) {
			continue
		}
		extra_files = append(extra_files, builder.NewStaticFileFromCore(core_types.StaticFile{
			Content:  jpeg_result.String(),
			Filename: fmt.Sprintf("%s.jpeg", shape.Nickname),
			Kind:     core_types.StaticFileUnknown,
		}))
	}
	build.AddStaticFiles(extra_files)

	return build
}

func SelectImageTga(shape *utfextract.Shape) (*utfextract.Image, error) {

	var result *utfextract.Image

	for _, image := range shape.Images {
		if image.Extension == "tga" {
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

	return nil, errors.New("not found image tga")
}
