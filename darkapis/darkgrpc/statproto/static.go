package statproto

import (
	_ "embed"

	"github.com/darklab8/fl-darkstat/darkcore/core_types"
)

//go:embed darkstat.swagger.json
var SwaggerContent string

var SwaggerJson core_types.StaticFile = core_types.StaticFile{
	Content:  SwaggerContent,
	Filename: "darkstat.swagger.json",
	Kind:     core_types.StaticFileUnknown,
}
