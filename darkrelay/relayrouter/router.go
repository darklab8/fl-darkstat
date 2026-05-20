package relayrouter

import (
	"github.com/darklab8/fl-darkstat/darkcore/builder"
	"github.com/darklab8/fl-darkstat/darkstat/appdata"
	"github.com/darklab8/fl-darkstat/darkstat/configs_export/infocarder"
	"github.com/darklab8/fl-darkstat/darkstat/front/tab"
	"github.com/darklab8/go-utils/utils/timeit"
	"github.com/darklab8/go-utils/utils/utils_types"
)

type Router struct {
	AppData *appdata.AppDataRelay
}

type RouterOpt func(l *Router)

func NewRouter(app_data *appdata.AppDataRelay, opts ...RouterOpt) *Router {
	l := &Router{AppData: app_data}
	for _, opt := range opts {
		opt(l)
	}

	return l
}

func (r *Router) Link() *builder.Builder {
	defer timeit.NewTimer("link, internal measure").Close()

	shared := r.AppData.Shared
	_ = shared
	configs := r.AppData.Configs
	build := appdata.NewBuilder(configs.Mapped.Discovery != nil)

	r.LinkPobs(r.AppData, build)

	timeit.NewTimerMF("linking most of stuff", func() {
		configs.GetInfocardsDict(func(infocards infocarder.Infocards) {
			for nickname, infocard := range infocards {
				build.RegComps(
					builder.NewComponent(
						utils_types.FilePath(tab.InfocardURL(nickname)),
						tab.Infocard(infocard),
					),
				)
			}
		})
	})

	return build
}
