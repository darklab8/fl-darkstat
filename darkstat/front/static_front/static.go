package static_front

import (
	_ "embed"
)

//go:embed html.min.js
var HtmxMinJs string

//go:embed sortable.js
var SortableJs string
