package darkmap

import (
	"context"
	"fmt"
	"os/signal"
	"syscall"

	"github.com/darklab8/fl-darkstat/darkcore/builder"
	"github.com/darklab8/fl-darkstat/darkcore/envers"
	"github.com/darklab8/fl-darkstat/darkcore/web"
	"github.com/darklab8/fl-darkstat/darkmap/front/urls"
	"github.com/darklab8/fl-darkstat/darkmap/linker"

	"github.com/darklab8/fl-darkstat/darkmap/settings"
	"github.com/darklab8/fl-darkstat/darkmap/settings/logus"

	"github.com/darklab8/go-utils/utils/cantil"
	"github.com/darklab8/go-utils/utils/timeit"
	"github.com/darklab8/go-utils/utils/utils_types"
)

func DarkmapCliGroup(Args []string) {
	urls.Index = utils_types.FilePath(settings.Env.IndexUrl)

	fmt.Println("freelancer folder=", settings.Env.FreelancerFolder, settings.Env)
	parser := cantil.NewConsoleParser(
		[]cantil.Action{
			{
				Nickname:    "build",
				Description: "build darkmap to static assets: html, css, js files",
				Func: func(info cantil.ActionInfo) error {
					linker.NewLinker(false).Link(context.Background()).BuildAll(false, nil)
					return nil
				},
			},
			{
				Nickname:    "web",
				Description: "run as standalone application that serves map from memory",
				Func: func(info cantil.ActionInfo) error {
					var fs *builder.Filesystem
					timer_web := timeit.NewTimer("total time for web web := func()")

					var linked_build *builder.Builder
					timer_NewLinkerLink := timeit.NewTimer("linking stuff linker.NewLinker().Link()")
					linked_build = linker.NewLinker(true).Link(context.Background())
					timer_NewLinkerLink.Close()

					timer_buildall := timeit.NewTimer("building stuff linked_build.BuildAll()")
					fs = linked_build.BuildAll(true, nil)
					timer_buildall.Close()

					timer_web.Close()
					graceful_closer := web.NewWeb(
						[]*builder.Filesystem{fs},
						web.WithSiteRoot(settings.Env.SiteRoot),
					).Serve(web.WebServeOpts{})

					ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
					defer stop()
					<-ctx.Done()

					graceful_closer.Close()

					return nil
				},
			},
		},
		cantil.ParserOpts{
			ParentArgs: []string{"darkmap"},
			Enverants:  envers.Enverants,
		},
	)
	err := parser.Run(Args)
	logus.Log.CheckError(err, "failed to execute darkmap cli group command")
}
