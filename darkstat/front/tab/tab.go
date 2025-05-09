package tab

import (
	"regexp"
	"strings"

	"github.com/darklab8/fl-darkstat/configs/configs_mapped/parserutils/inireader"
	"github.com/darklab8/fl-darkstat/darkstat/configs_export/infocarder"
)

type ShowEmpty bool

func InfocardURL(infocard_key infocarder.InfocardKey) string {
	return "cdn/infocards/info_" + strings.ToLower(string(infocard_key))
}

func GetFirstLine(infocards *infocarder.Infocarder, infokey infocarder.InfocardKey) string {
	if infocard_lines, ok := infocards.GetInfocard2(infokey); ok {
		if len(infocard_lines) > 0 {
			return string(infocard_lines[0].ToStr())
		}
	}
	return ""
}

func GetShipName(infocards *infocarder.Infocarder, infokey infocarder.InfocardKey) string {
	first_line := GetFirstLine(infocards, infokey)
	var result string

	if found := ShipNameRegex.FindStringSubmatch(first_line); len(found) > 0 {
		result = found[1]
	}
	return result
}

var ShipNameRegex *regexp.Regexp

func init() {
	inireader.InitRegexExpression(&ShipNameRegex, `\"(.*)\"`)
}
