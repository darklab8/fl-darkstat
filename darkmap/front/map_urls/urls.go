package map_urls

import (
	"strings"

	"github.com/darklab8/go-utils/utils/utils_types"
)

var (
	Index     utils_types.FilePath = "index.html"
	SearchBar utils_types.FilePath = "search-bar"
)

func SystemDetailedUrlNick(system_nick string) string {
	return "cdn/map/system/system-" + strings.ToLower(system_nick) + ".html"
}

func ObjUrl(system_nick string, obj_nickname string) string {
	return SystemDetailedUrlNick(system_nick) + "?q=" + obj_nickname
}
