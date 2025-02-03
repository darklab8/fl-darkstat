package tab

import (
	"strings"

	"github.com/darklab8/fl-darkstat/darkstat/configs_export"
)

type ShowEmpty bool

func InfocardURL(infocard_key configs_export.InfocardKey) string {
	return "infocards/info_" + strings.ToLower(string(infocard_key))
}

func GetFirstLine(infocards configs_export.Infocards, infokey configs_export.InfocardKey) string {
	if infocard_lines, ok := infocards[infokey]; ok {
		if len(infocard_lines) > 0 {
			return string(infocard_lines[0].ToStr())
		}
	}
	return ""
}
