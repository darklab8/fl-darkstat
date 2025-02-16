package relayrouter

import (
	"github.com/darklab8/fl-darkstat/darkcore/builder"
	"github.com/darklab8/fl-darkstat/darkrelay/relayfront"
	"github.com/darklab8/fl-darkstat/darkstat/appdata"
	"github.com/darklab8/fl-darkstat/darkstat/front/tab"
	"github.com/darklab8/fl-darkstat/darkstat/front/types"
	"github.com/darklab8/fl-darkstat/darkstat/front/urls"
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
	build := r.AppData.Build
	configs := r.AppData.Configs

	r.LinkPobs(r.AppData)

	build.RegComps(
		builder.NewComponent(
			urls.Index,
			relayfront.Index(types.ThemeLight, shared),
		),
		builder.NewComponent(
			urls.DarkIndex,
			relayfront.Index(types.ThemeDark, shared),
		),
		builder.NewComponent(
			urls.VanillaIndex,
			relayfront.Index(types.ThemeVanilla, shared),
		),
	)

	timeit.NewTimerMF("linking most of stuff", func() {
		for nickname, infocard := range configs.Infocards {
			build.RegComps(
				builder.NewComponent(
					utils_types.FilePath(tab.InfocardURL(nickname)),
					tab.Infocard(infocard),
				),
			)
		}
	})

	return build
}
