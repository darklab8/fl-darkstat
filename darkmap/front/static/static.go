package static

import (
	"embed"

	"github.com/darklab8/fl-darkstat/darkcore/core_front"
	"github.com/darklab8/go-utils/utils/utils_types"
)

//go:embed files/*
var currentdir embed.FS

var StaticFilesystem core_front.StaticFilesystem = core_front.GetFiles(
	currentdir,
	utils_types.GetFilesParams{
		RootFolder:        utils_types.FilePath("files"),
		AllowedExtensions: []string{"js", "css", "png", "jpeg", "jpg"},
	},
)
