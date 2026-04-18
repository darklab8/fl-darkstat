package linker

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"html/template"
	"sort"
	"strings"
	"time"

	"github.com/darklab8/fl-darkstat/darkcore/builder"
	"github.com/darklab8/fl-darkstat/darkcore/core_static"
	"github.com/darklab8/fl-darkstat/darkcore/core_types"
	"github.com/darklab8/fl-darkstat/darkstat/configs_export/infocarder"

	"github.com/darklab8/fl-darkstat/darkmap/export_map"
	"github.com/darklab8/fl-darkstat/darkmap/front"
	"github.com/darklab8/fl-darkstat/darkmap/front/map_urls"
	"github.com/darklab8/fl-darkstat/darkmap/front/static"
	"github.com/darklab8/fl-darkstat/darkmap/front/static_front"
	"github.com/darklab8/fl-darkstat/darkmap/search_bar"
	"github.com/darklab8/fl-darkstat/darkmap/settings"
	"github.com/darklab8/fl-darkstat/darkmap/settings/logus"
	"github.com/darklab8/fl-darkstat/darkmap/types"
	"github.com/darklab8/fl-darkstat/darkmap/utfextract"
	stat_front "github.com/darklab8/fl-darkstat/darkstat/front/static_front"
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
	if l.Export == nil {
		l.Export = export_map.NewExport(ctx)
	}

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

	tmpl, err := template.New("zones").Parse(static_front.ZonesCSS.Content)
	logus.Log.CheckPanic(err, "failed to parse zones css")
	var templated_zones bytes.Buffer
	err = tmpl.Execute(&templated_zones, params)
	logus.Log.CheckPanic(err, "failed to template zones css")

	files := []builder.StaticFile{
		builder.NewStaticFileFromCore(core_static.HtmxJS),
		builder.NewStaticFileFromCore(core_static.HtmxPreloadJS),
		builder.NewStaticFileFromCore(core_static.SortableJS),
		builder.NewStaticFileFromCore(core_static.ResetCSS),
		builder.NewStaticFileFromCore(static_front.FaviconIco),
		builder.NewStaticFileFromCore(static_front.CommonCSS),
		builder.NewStaticFileFromCore(static_front.CustomCSS),
		builder.NewStaticFileFromCore(stat_front.CustomCrossCSS),
		builder.NewStaticFileFromCore(static_front.GalaxyCSS),
		builder.NewStaticFileFromCore(static_front.CustomJS),
		builder.NewStaticFileFromCore(static_front.MapGalaxyJS),
		builder.NewStaticFileFromCore(static_front.MapSystemJS),
		builder.NewStaticFileFromCore(static_front.PanzoomJS),
		builder.NewStaticFileFromCore(static_front.RemodalCSS),
		builder.NewStaticFileFromCore(search_bar.SearchBarCss),
		builder.NewStaticFileFromCore(stat_front.CustomJSSCrossShared),
		builder.NewStaticFileFromCore(static_front.ZonesCSS.GetTemplated(templated_zones.String())),
		builder.NewStaticFileFromCore(static_front.TippyJS1),
		builder.NewStaticFileFromCore(static_front.TippyJS2),
		builder.NewStaticFileFromCore(core_types.StaticFile{
			Content:  core_static.FaviconIcoContent,
			Filename: "stat_favicon.ico",
			Kind:     core_types.StaticFileIco,
		}),
	}

	for _, file := range static.StaticFilesystem.Files {
		files = append(files, builder.NewStaticFileFromCore(file))
	}
	build = builder.NewBuilder(params, files)

	map_search_entries := l.Export.SearchEntries
	var search_entries []search_bar.Entry

	for _, value := range map_search_entries {
		search_entries = append(search_entries, value)
	}
	sort.Slice(search_entries, func(i, j int) bool {
		if search_entries[i].Kind != search_entries[j].Kind {
			return search_entries[i].Name < search_entries[j].Name
		}
		return search_entries[i].Name < search_entries[j].Name
	})

	build.RegComps(
		builder.NewComponent(
			map_urls.Index,
			front.Index(l.Export),
		),
		builder.NewComponent(
			map_urls.SearchBar,
			search_bar.SearchBar(search_entries, l.Export.Mapped.Discovery != nil),
		),
	)

	for _, system := range l.Export.Systems {
		build.RegComps(
			builder.NewComponent(
				utils_types.FilePath(front.SystemDetailedUrl(system)),
				front.System(system, l.Export),
			),
		)
	}

	timeit.NewTimerMF("linking most of stuff", func() {
		l.Export.Exp.GetInfocardsDict(func(infocards infocarder.Infocards) {
			for nickname, infocard := range infocards {
				build.RegComps(
					builder.NewComponent(
						utils_types.FilePath(front.MapInfocardURL(nickname)),
						front.MapInfocard(infocard),
					),
				)
			}
		})
	})

	build.RegComps(
		builder.NewComponent(
			utils_types.FilePath(front.MapInfocardURL("map_legend")),
			front.MapLegend(l.Export),
		),
	)

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
		// if l.ToMemory {
		// 	if _, permitted := l.Export.Shapes.PermittedShapes[strings.ToLower(shape.Nickname)]; !permitted {
		// 		continue
		// 	}
		// }
		if _, permitted := l.Export.Shapes.PermittedShapes[strings.ToLower(shape.Nickname)]; !permitted {
			continue
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

			logus.Log.Info("image linking", typelog.Any("dest", image.Dest), typelog.Any("nick", image.Nickname))

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

	var failed_files []StaticFileInParallel

	for i := 0; i < created_jobs; i++ {
		result := <-decoded_shape_files
		if !result.made {
			failed_files = append(failed_files, result)
			continue
		}
		extra_files = append(extra_files, builder.NewStaticFileFromCore(result.StaticFile))
	}

	fmt.Println("SHAPES FINISHED ACCEPTING ALL JOBS, took time seconds=", time.Since(time_start).Seconds(), " handled jobs=", created_jobs)

	build.AddStaticFiles(extra_files)

	type FallBackInfo struct {
		Count int
		Info  string
	}

	found_fallbacks := make(map[string]FallBackInfo)
	not_found_fallbacks := make(map[string]FallBackInfo)
	file_checker := build.GetStaticFileChecker()
	for _, system := range l.Export.Systems {
		for _, obj := range system.Objs {
			if _, ok := file_checker[utils_types.FilePath(fmt.Sprintf("%s.png", utils_types.FilePath(obj.ShapeName)))]; !ok {
				if obj.Kind == export_map.ObjPlanet {

					switch obj.ShapeName {
					// optionally use images as fallbacks
					// // DISCO START
					case "ast_lava_hd":
						obj.ShapeName = "fallback/ast_lava"
					case "ast_rock":
						obj.ShapeName = "fallback/ast_rock"
					case "solar_mat_dyson_city":
						obj.UseFallback = true
					case "detailmap_crater_mid":
						obj.ShapeName = "fallback/detailmap_crater_mid"
					// // DISCO END

					default:
						obj.UseFallback = true
						var count int
						if value, ok := not_found_fallbacks[obj.ShapeName]; ok {
							count = value.Count
						}
						not_found_fallbacks[obj.ShapeName] = FallBackInfo{Count: count + 1, Info: fmt.Sprint(strings.Join([]string{"planet", system.Name, obj.Nickname}, ","))}
					}

				} else if obj.Kind != export_map.ObjPlanet {

					switch obj.ShapeName {
					// optionally use images as fallbacks

					// // VANILLA START
					case "nnm_sm_mining":
						obj.UseFallback = true
					case "nnm_sm_depot":
						obj.ShapeName = "fallback/nav_depot"
					case "nnm_sm_navbuoy": // colored by css
						obj.UseFallback = true
					case "nav_buoy": // colored by css
						obj.UseFallback = true
					case "nnm_sm_communications":
						obj.ShapeName = "fallback/nav_lootabledepot"
					case "nnm_sm_mplatform":
						obj.ShapeName = "fallback/nav_lootabledepot"
					// // VANILLA END

					default:
						obj.UseFallback = true
						var count int
						if value, ok := not_found_fallbacks[obj.ShapeName]; ok {
							count = value.Count
						}

						not_found_fallbacks[obj.ShapeName] = FallBackInfo{Count: count + 1, Info: fmt.Sprint(strings.Join([]string{"obj", system.Name, system.Nickname, obj.Nickname}, ","))}
					}

				}
			}
		}
	}

	for key, value := range found_fallbacks {
		logus.Log.Warn("found static file fallback for object", typelog.Any("key", key), typelog.Any("value", value))
	}
	for key, value := range not_found_fallbacks {
		logus.Log.Warn("not found static file fallback for object", typelog.Any("key", key), typelog.Any("value", value))
	}

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
