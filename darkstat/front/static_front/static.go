package static_front

import (
	_ "embed"

	"github.com/darklab8/fl-darkstat/darkcore/core_types"
)

//go:embed custom/shared_vanilla.js
var CustomJSCSharedVanilla string

var CustomJSSharedVanilla core_types.StaticFile = core_types.StaticFile{
	Content:  CustomJSCSharedVanilla,
	Filename: "custom/shared_vanilla.js",
	Kind:     core_types.StaticFileJS,
}

//go:embed custom/shared_discovery.js
var CustomJSCSharedDiscovery string

var CustomJSSharedDiscovery core_types.StaticFile = core_types.StaticFile{
	Content:  CustomJSCSharedDiscovery,
	Filename: "custom/shared_discovery.js",
	Kind:     core_types.StaticFileJS,
}

//go:embed custom/shared.js
var CustomJSCShared string

var CustomJSShared core_types.StaticFile = core_types.StaticFile{
	Content:  CustomJSCShared,
	Filename: "custom/shared.js",
	Kind:     core_types.StaticFileJS,
}

//go:embed custom/main.js
var CustomJSContent string

var CustomJS core_types.StaticFile = core_types.StaticFile{
	Content:  CustomJSContent,
	Filename: "custom/main.js",
	Kind:     core_types.StaticFileJS,
}

//go:embed custom/table_resizer.js
var CustomResizerJSContent string

var CustomJSResizer core_types.StaticFile = core_types.StaticFile{
	Content:  CustomResizerJSContent,
	Filename: "table_resizer.js",
	Kind:     core_types.StaticFileJS,
}

//go:embed custom/filtering.js
var CustomFilteringJS string

var CustomJSFiltering core_types.StaticFile = core_types.StaticFile{
	Content:  CustomFilteringJS,
	Filename: "filtering.js",
	Kind:     core_types.StaticFileJS,
}

//go:embed custom/filter_route_min_dists.js
var CustomFilteringRoutesJS string

var CustomJSFilteringRoutes core_types.StaticFile = core_types.StaticFile{
	Content:  CustomFilteringRoutesJS,
	Filename: "filter_route_min_dists.js",
	Kind:     core_types.StaticFileJS,
}

//go:embed common.css
var CommonCSSContent string

var CommonCSS core_types.StaticFile = core_types.StaticFile{
	Content:  CommonCSSContent,
	Filename: "common.css",
	Kind:     core_types.StaticFileCSS,
}

//go:embed custom.css
var CustomCSSContent string

var CustomCSS core_types.StaticFile = core_types.StaticFile{
	Content:  CustomCSSContent,
	Filename: "custom.css",
	Kind:     core_types.StaticFileCSS,
}
