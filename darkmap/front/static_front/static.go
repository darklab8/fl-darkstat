package static_front

import (
	_ "embed"

	"github.com/darklab8/fl-darkstat/darkcore/core_types"
)

//go:embed panzoom.min.js
var PanzoomJSContent string

var PanzoomJS core_types.StaticFile = core_types.StaticFile{
	Content:  PanzoomJSContent,
	Filename: "panzoom.min.js",
	Kind:     core_types.StaticFileJS,
}

//go:embed custom.js
var CustomJSContent string

var CustomJS core_types.StaticFile = core_types.StaticFile{
	Content:  CustomJSContent,
	Filename: "custom.js",
	Kind:     core_types.StaticFileJS,
}

//go:embed map_galaxy.js
var MapGalaxyJSContent string

var MapGalaxyJS core_types.StaticFile = core_types.StaticFile{
	Content:  MapGalaxyJSContent,
	Filename: "map_galaxy.js",
	Kind:     core_types.StaticFileJS,
}

//go:embed map_system.js
var MapSystemJSContent string

var MapSystemJS core_types.StaticFile = core_types.StaticFile{
	Content:  MapSystemJSContent,
	Filename: "map_system.js",
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

//go:embed galaxy.css
var GalaxyCSSContent string

var GalaxyCSS core_types.StaticFile = core_types.StaticFile{
	Content:  GalaxyCSSContent,
	Filename: "galaxy.css",
	Kind:     core_types.StaticFileCSS,
}

//go:embed favicon.ico
var FaviconIcoContent string

var FaviconIco core_types.StaticFile = core_types.StaticFile{
	Content:  FaviconIcoContent,
	Filename: "favicon.ico",
	Kind:     core_types.StaticFileIco,
}
