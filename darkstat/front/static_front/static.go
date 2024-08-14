package static_front

import (
	_ "embed"

	"github.com/darklab8/fl-darkcore/darkcore/core_types"
)

//go:embed custom.js
var CustomJSContent string

var CustomJS core_types.StaticFile = core_types.StaticFile{
	Content:  CustomJSContent,
	Filename: "custom.js",
	Kind:     core_types.StaticFileJS,
}

//go:embed common.html
var CommonCSSContent string

var CommonCSS core_types.StaticFile = core_types.StaticFile{
	Content:  CommonCSSContent,
	Filename: "common.html",
	Kind:     core_types.StaticFileCSS,
}

//go:embed custom.html
var CustomCSSContent string

var CustomCSS core_types.StaticFile = core_types.StaticFile{
	Content:  CustomCSSContent,
	Filename: "custom.html",
	Kind:     core_types.StaticFileCSS,
}
