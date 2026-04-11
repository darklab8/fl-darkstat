package urls

import (
	"strings"

	"github.com/darklab8/go-utils/utils/utils_types"
)

var (
	Index utils_types.FilePath = "index.html"
)

func SystemDetailedUrlNick(system_nick string) string {
	return "cdn/map/system/system-" + strings.ToLower(system_nick) + ".html"
}
