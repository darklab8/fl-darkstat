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

func themeCycleOrder(priority Theme) [3]Theme {
	switch priority {
	case ThemeDark:
		return [3]Theme{ThemeDark, ThemeVanilla, ThemeLight}
	case ThemeVanilla:
		return [3]Theme{ThemeVanilla, ThemeLight, ThemeDark}
	default:
		return [3]Theme{ThemeLight, ThemeDark, ThemeVanilla}
	}
}

func ThemeCycleURLs(siteRoot string, priority Theme) []string {
	order := themeCycleOrder(priority)
	return []string{
		siteRoot + ThemeIndexHTMLFile(order[0]),
		siteRoot + ThemeIndexHTMLFile(order[1]),
		siteRoot + ThemeIndexHTMLFile(order[2]),
	}
}

func ThemeCycleNicks(priority Theme) []string {
	order := themeCycleOrder(priority)
	return []string{order[0].ToNick(), order[1].ToNick(), order[2].ToNick()}
}

func ThemeStorageNicks() []string {
	return []string{ThemeLight.ToNick(), ThemeDark.ToNick(), ThemeVanilla.ToNick()}
}

func ThemeStorageIndexFiles() []string {
	return []string{ThemeIndexHTMLFile(ThemeLight), ThemeIndexHTMLFile(ThemeDark), ThemeIndexHTMLFile(ThemeVanilla)}
}
