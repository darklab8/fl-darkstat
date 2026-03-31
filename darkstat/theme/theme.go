package theme

import (
	"fmt"
	"strings"
)

type Theme int64

const (
	ThemeNotSet Theme = iota
	ThemeLight
	ThemeDark
	ThemeVanilla
)

func (t Theme) ToNick() string {
	switch t {
	case ThemeLight:
		return "light"
	case ThemeDark:
		return "dark"
	case ThemeVanilla:
		return "vanilla"
	default:
		panic(fmt.Sprintf("not a valid theme: %d", int64(t)))
	}
}

func ThemeIndexHTMLFile(t Theme) string {
	return t.ToNick() + ".html"
}

func ParseDefaultThemeName(s string) Theme {
	switch strings.ToLower(strings.TrimSpace(s)) {
	case "white", "light":
		return ThemeLight
	case "dark":
		return ThemeDark
	case "vanilla":
		return ThemeVanilla
	default:
		panic(fmt.Sprintf("unrecognized default theme name: %q", s))
	}
}
