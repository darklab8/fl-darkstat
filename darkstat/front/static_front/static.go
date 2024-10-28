package static_front

import (
	_ "embed"

	"github.com/darklab8/fl-darkcore/darkcore/core_types"
)

//go:embed custom/shared_vanilla.js
var CustomJSCSharedVanilla string

var CustomJSSharedVanilla core_types.StaticFile = core_types.StaticFile{
	Content:  CustomJSCSharedVanilla,
	Filename: "custom_shared_vanilla.js",
	Kind:     core_types.StaticFileJS,
}

//go:embed custom/shared_discovery.js
var CustomJSCSharedDiscovery string

var CustomJSSharedDiscovery core_types.StaticFile = core_types.StaticFile{
	Content:  CustomJSCSharedDiscovery,
	Filename: "custom_shared_discovery.js",
	Kind:     core_types.StaticFileJS,
}

//go:embed custom/shared.js
var CustomJSCShared string

var CustomJSShared core_types.StaticFile = core_types.StaticFile{
	Content:  CustomJSCShared,
	Filename: "custom_shared.js",
	Kind:     core_types.StaticFileJS,
}

//go:embed custom/main.js
var CustomJSContent string

var CustomJS core_types.StaticFile = core_types.StaticFile{
	Content:  CustomJSContent,
	Filename: "custom_main.js",
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

//go:embed docs/docs_tech_compat_id_selector.png
var PictureTechCompatIDSelector string

//go:embed docs/docs_techcompat_hover.png
var PictureTechCompatHover string

//go:embed docs/docs_coordinates_in_trade_routes.png
var PictureCoordinatesTradeRoutes string

//go:embed docs/docs_timestamp.png
var Picture4 string

//go:embed docs/docs_pinning.png
var Picture5 string

//go:embed docs/docs_movable_borders.png
var Picture6 string

//go:embed docs/docs_search_bar.png
var Picture7 string

//go:embed docs/docs_infocard_search.png
var Picture8 string

var Pictures []core_types.StaticFile = []core_types.StaticFile{
	{
		Content:  PictureTechCompatIDSelector,
		Filename: "docs_tech_compat_id_selector.png",
		Kind:     core_types.StaticFilePicture,
	},
	{
		Content:  PictureTechCompatHover,
		Filename: "docs_techcompat_hover.png",
		Kind:     core_types.StaticFilePicture,
	},
	{
		Content:  PictureCoordinatesTradeRoutes,
		Filename: "docs_coordinates_in_trade_routes.png",
		Kind:     core_types.StaticFilePicture,
	},
	{
		Content:  Picture4,
		Filename: "docs_timestamp.png",
		Kind:     core_types.StaticFilePicture,
	},
	{
		Content:  Picture5,
		Filename: "docs_pinning.png",
		Kind:     core_types.StaticFilePicture,
	},
	{
		Content:  Picture6,
		Filename: "docs_movable_borders.png",
		Kind:     core_types.StaticFilePicture,
	},
	{
		Content:  Picture7,
		Filename: "docs_search_bar.png",
		Kind:     core_types.StaticFilePicture,
	},
	{
		Content:  Picture8,
		Filename: "docs_infocard_search.png",
		Kind:     core_types.StaticFilePicture,
	},
}
