package relayrouter

import (
	"github.com/darklab8/fl-darkcore/darkcore/builder"
	"github.com/darklab8/fl-darkstat/darkrelay/relayfront"
	"github.com/darklab8/fl-darkstat/darkstat/front/types"
	"github.com/darklab8/fl-darkstat/darkstat/front/urls"
	"github.com/darklab8/fl-darkstat/darkstat/router"
	"github.com/darklab8/go-utils/utils/timeit"
)

type Router struct {
	AppData *router.AppData
}

type RouterOpt func(l *Router)

func NewRouter(app_data *router.AppData, opts ...RouterOpt) *Router {
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

	return build
}
