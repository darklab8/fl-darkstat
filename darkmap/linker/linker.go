package linker

import (
	"time"

	"github.com/darklab8/fl-darkstat/darkcore/builder"
	"github.com/darklab8/fl-darkstat/darkcore/core_static"

	"github.com/darklab8/fl-darkstat/darkmap/front"
	"github.com/darklab8/fl-darkstat/darkmap/front/export_front"
	"github.com/darklab8/fl-darkstat/darkmap/front/static"
	"github.com/darklab8/fl-darkstat/darkmap/front/static_front"
	"github.com/darklab8/fl-darkstat/darkmap/front/urls"
	"github.com/darklab8/fl-darkstat/darkmap/settings"
	"github.com/darklab8/fl-darkstat/darkmap/types"
	"github.com/darklab8/go-utils/utils/timeit"
)

type Linker struct {
	Export *export_front.Export
}

type LinkOption func(l *Linker)

func NewLinker(opts ...LinkOption) *Linker {
	l := &Linker{}
	for _, opt := range opts {
		opt(l)
	}

	return l
}

func (l *Linker) Link() *builder.Builder {
	l.Export = export_front.NewExport()

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
		builder.NewStaticFileFromCore(core_static.FaviconIco),
		builder.NewStaticFileFromCore(static_front.CommonCSS),
		builder.NewStaticFileFromCore(static_front.CustomCSS),
		builder.NewStaticFileFromCore(static_front.CustomJS),
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

	return build
}
