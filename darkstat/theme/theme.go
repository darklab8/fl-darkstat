package theme

import (
	"fmt"
	"os"
	"strings"
)

type Theme int64

const (
	ThemeNotSet Theme = iota
	ThemeLight
	ThemeDarkLight
	ThemeDark
	ThemeVanilla
)

func (t Theme) ToNick() string {
	switch t {
	case ThemeLight:
		return "light"
	case ThemeDark:
		return "dark"
	case ThemeDarkLight:
		return "darklight"
	case ThemeVanilla:
		return "vanilla"
	default:
		panic(fmt.Sprintf("not a valid theme: %d", int64(t)))
	}
}

func ThemeIndexHTMLFile(t Theme) string {
	if value, ok := os.LookupEnv("theme_" + t.ToNick() + "_url"); ok {
		return value
	}
	return t.ToNick() + ".html"
}

func ParseDefaultThemeName(s string) Theme {
	switch strings.ToLower(strings.TrimSpace(s)) {
	case "white", "light":
		return ThemeLight
	case "dark":
		return ThemeDark
	case "darklight":
		return ThemeDarkLight
	case "vanilla":
		return ThemeVanilla
	default:
		panic(fmt.Sprintf("unrecognized default theme name: %q", s))
	}
}

func ThemeCycleOrder() (result []Theme) {
	result = []Theme{ThemeLight, ThemeDark, ThemeDarkLight, ThemeVanilla}
	return
}
