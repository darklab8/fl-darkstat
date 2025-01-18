package tab

import (
	"strings"

	"github.com/darklab8/fl-darkstat/darkstat/configs_export"
)

type ShowEmpty bool

func InfocardURL(infocard_key configs_export.InfocardKey) string {
	return "infocards/info_" + strings.ToLower(string(infocard_key))
}
