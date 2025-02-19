package staticproto

import (
	_ "embed"
)

//go:embed index.html
var Index string

//go:embed swagger-ui-bundle.js
var JS1 string

//go:embed swagger-ui-standalone-preset.js
var JS2 string

//go:embed swagger-ui.css
var CSS string
