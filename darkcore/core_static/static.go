package core_static

import (
	_ "embed"

	"github.com/darklab8/fl-darkstat/darkcore/core_types"
)

// Commented out IE stuff as it makes things slow
// if (element.children) { // IE
//
//		forEach(element.children, function(child) { cleanUpElement(child) });
//	}
//
// see https://github.com/bigskysoftware/htmx/issues/879 for more details
//
// also commented out  //   handleAttributes(parentNode, fragment, settleInfo);
// because we don't need CSS transitions and they are hurtful https://htmx.org/docs/#css_transitions
//

//go:embed htmx.1.9.11.js
var HtmxJsContent string

var HtmxJS core_types.StaticFile = core_types.StaticFile{
	Content:  HtmxJsContent,
	Filename: "htmx.js",
	Kind:     core_types.StaticFileJS,
}

//go:embed htmx.1.9.11.preload.js
var PreloadJsContent string

var HtmxPreloadJS core_types.StaticFile = core_types.StaticFile{
	Content:  PreloadJsContent,
	Filename: "htmx_preload.js",
	Kind:     core_types.StaticFileJS,
}

//go:embed sortable.js
var SortableJsContent string

var SortableJS core_types.StaticFile = core_types.StaticFile{
	Content:  SortableJsContent,
	Filename: "sortable.js",
	Kind:     core_types.StaticFileJS,
}

//go:embed reset.css
var ResetCSSContent string

var ResetCSS core_types.StaticFile = core_types.StaticFile{
	Content:  ResetCSSContent,
	Filename: "reset.css",
	Kind:     core_types.StaticFileCSS,
}

//go:embed favicon.ico
var FaviconIcoContent string

var FaviconIco core_types.StaticFile = core_types.StaticFile{
	Content:  FaviconIcoContent,
	Filename: "favicon.ico",
	Kind:     core_types.StaticFileIco,
}
