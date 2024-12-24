package router

import (
	"github.com/a-h/templ"
	"github.com/darklab8/fl-darkcore/darkcore/builder"
	"github.com/darklab8/fl-darkstat/darkstat/front/tab"
	"github.com/darklab8/go-utils/utils/utils_types"
)

type TabRouter[T any] struct {
	build          *builder.Builder
	items          []T
	filtered_items []T
}

type TabRouterOpt[T any] func(l *TabRouter[T])

type TabModeI interface {
	ToInt() int64
}

func NewTabRouter[T any](build *builder.Builder, items []T, filter_to_useful func(items []T) []T, opts ...TabRouterOpt[T]) *TabRouter[T] {
	tab_router := &TabRouter[T]{
		build:          build,
		items:          items,
		filtered_items: filter_to_useful(items),
	}

	for _, opt := range opts {
		opt(tab_router)
	}

	return tab_router
}

func (t *TabRouter[T]) Register(
	url utils_types.FilePath,
	callback func(items []T, show_empty tab.ShowEmpty) templ.Component,
) {

	t.build.RegComps(
		builder.NewComponent(
			url,
			callback(t.filtered_items, tab.ShowEmpty(false)),
		),
		builder.NewComponent(
			tab.AllItemsUrl(url),
			callback(t.items, tab.ShowEmpty(true)),
		),
	)
}
