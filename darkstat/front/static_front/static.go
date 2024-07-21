package static_front

import (
	_ "embed"
)

//go:embed custom.js
var CustomJS string

//go:embed favicon.ico
var FaviconIco string
