package static_front

import (
	_ "embed"
)

//go:embed htmx.1.9.11.min.js
var HtmxMinJs string

//go:embed htmx.1.9.11.preload.js
var PreloadJs string

//go:embed sortable.js
var SortableJs string

//go:embed custom.js
var CustomJS string
