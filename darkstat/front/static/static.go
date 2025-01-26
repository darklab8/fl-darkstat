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
	utils_types.GetFilesParams{RootFolder: utils_types.FilePath("files")},
)

// Example how to import on init
// var PictureCoordinatesInTradeRoutes core_types.StaticFile = StaticFilesystem.GetFileByRelPath("docs_coordinates_in_trade_routes.png")
