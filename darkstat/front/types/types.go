package types

import (
	"context"
	"strings"
	"time"
	"unicode"

	"github.com/darklab8/fl-darkstat/configs/cfg"
	"github.com/darklab8/fl-darkstat/configs/discovery/minecontrol"
	"github.com/darklab8/fl-darkstat/configs/discovery/techcompat"
	"github.com/darklab8/fl-darkstat/darkcore/core_types"
	"github.com/darklab8/fl-darkstat/darkstat/configs_export"
	"github.com/darklab8/fl-darkstat/darkstat/configs_export/infocarder"
	"github.com/darklab8/fl-data-discovery/autopatcher"
	"github.com/darklab8/go-utils/utils/utils_types"
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
		return ""
	}
}

func ThemeIndexHTMLFile(t Theme) string {
	if n := t.ToNick(); n != "" {
		return n + ".html"
	}
	return "light.html"
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
		return ThemeLight
	}
}

func ThemeCycleURLs(siteRoot string, priority Theme) []string {
	var order [3]Theme
	switch priority {
	case ThemeDark:
		order = [3]Theme{ThemeDark, ThemeVanilla, ThemeLight}
	case ThemeVanilla:
		order = [3]Theme{ThemeVanilla, ThemeLight, ThemeDark}
	default:
		order = [3]Theme{ThemeLight, ThemeDark, ThemeVanilla}
	}
	return []string{
		siteRoot + ThemeIndexHTMLFile(order[0]),
		siteRoot + ThemeIndexHTMLFile(order[1]),
		siteRoot + ThemeIndexHTMLFile(order[2]),
	}
}

type GlobalParams struct {
	Buildpath      utils_types.FilePath
	Theme          Theme
	Themes         []string
	SiteHost       string
	SiteRoot       string
	SiteUrl        string
	StaticRoot     string
	Heading        string
	Timestamp      time.Time
	TractorTabName string

	RelayHost string
	RelayRoot string

	AppStart time.Time

	ShowDisco bool
}

func (g *GlobalParams) GetBuildPath() utils_types.FilePath {
	return g.Buildpath
}

func (g *GlobalParams) GetStaticRoot() string {
	return g.StaticRoot
}

func GetCtx(ctx context.Context) *GlobalParams {
	return ctx.Value(core_types.GlobalParamsCtxKey).(*GlobalParams)
}

type FLSRData struct {
	ShowFLSR bool
}

type DiscoveryData struct {
	ShowDisco    bool
	Ids          []*configs_export.Tractor
	TractorsByID map[cfg.TractorID]*configs_export.Tractor
	Config       *techcompat.Config
	LatestPatch  autopatcher.Patch
	Minecontrol  *minecontrol.Config

	*infocarder.Infocarder

	OrderedTechcompat configs_export.TechCompatOrderer
}

type ShipNames struct {
	Transport string
	Frigate   string
	Freighter string
}
type SharedData struct {
	DiscoveryData
	FLSRData
	CraftableBaseName     string
	AverageTradeLaneSpeed int
	ShipNames             ShipNames
}

// COPY PASTE FROM https://cs.opensource.google/go/go/+/refs/tags/go1.21.10:src/strings/strings.go;l=752
// isSeparator reports whether the rune could mark a word boundary.
// TODO: update when package unicode captures more of the properties.
func isSeparator(r rune) bool {
	// ASCII alphanumerics and underscore are not separators
	if r <= 0x7F {
		switch {
		case '0' <= r && r <= '9':
			return false
		case 'a' <= r && r <= 'z':
			return false
		case 'A' <= r && r <= 'Z':
			return false
		case r == '_':
			return false
		}
		return true
	}
	// Letters and digits are not separators
	if unicode.IsLetter(r) || unicode.IsDigit(r) {
		return false
	}
	// Otherwise, all we can do for now is treat spaces as separators.
	return unicode.IsSpace(r)
}

// COPY PASTE FROM https://cs.opensource.google/go/go/+/refs/tags/go1.21.10:src/strings/strings.go;l=782
func Title(s string) string {
	// Use a closure here to remember state.
	// Hackish but effective. Depends on Map scanning in order and calling
	// the closure once per rune.
	prev := ' '
	return strings.Map(
		func(r rune) rune {
			if isSeparator(prev) {
				prev = r
				return unicode.ToTitle(r)
			}
			prev = r
			return r
		},
		s)
}

func ToCapital(value string) string {
	return Title(value)
}
