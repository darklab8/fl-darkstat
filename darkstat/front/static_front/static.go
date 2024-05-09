package static_front

import (
	_ "embed"
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
var HtmxJs string

//go:embed htmx.1.9.11.preload.js
var PreloadJs string

//go:embed sortable.js
var SortableJs string

//go:embed custom.js
var CustomJS string
