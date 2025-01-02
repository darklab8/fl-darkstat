package relayrouter

import (
	"github.com/darklab8/fl-darkcore/darkcore/builder"
	"github.com/darklab8/fl-darkstat/darkrelay/relayfront"
	"github.com/darklab8/fl-darkstat/darkstat/front/types"
	"github.com/darklab8/fl-darkstat/darkstat/front/urls"
	"github.com/darklab8/fl-darkstat/darkstat/router"
	"github.com/darklab8/go-utils/utils/ptr"
	"github.com/darklab8/go-utils/utils/timeit"
)

type Router struct {
	AppData *router.AppData
}

type RouterOpt func(l *Router)

func NewRouter(opts ...RouterOpt) *Router {
	l := &Router{}
	for _, opt := range opts {
		opt(l)
	}

	return l
}

func WithAppData(AppData *router.AppData) RouterOpt {
	return func(l *Router) { l.AppData = AppData }
}

func (r *Router) Link() *builder.Builder {
	defer timeit.NewTimer("link, internal measure").Close()

	if r.AppData == nil {
		r.AppData = router.NewAppData()
	}

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

func (r *Router) Update() {
	r.AppData.Configs.PoBs = r.AppData.Configs.GetPoBs()

	for _, pob := range r.AppData.Configs.PoBs {
		if pob.Money == nil {
			pob.Money = ptr.Ptr(1)
		}
		pob.Money = ptr.Ptr(*pob.Money + 1)
	}
}
