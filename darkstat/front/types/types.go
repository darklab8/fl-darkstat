package types

import (
	"context"
	"time"

	"github.com/darklab8/fl-darkstat/configs/cfg"
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

type GlobalParams struct {
	Buildpath      utils_types.FilePath
	Theme          Theme
	Themes         []string
	SiteHost       string
	SiteRoot       string
	StaticRoot     string
	Heading        string
	Timestamp      time.Time
	TractorTabName string

	RelayHost string
	RelayRoot string
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

	*infocarder.Infocarder

	OrderedTechcompat configs_export.TechCompatOrderer
}

type SharedData struct {
	DiscoveryData
	FLSRData
	CraftableBaseName     string
	AverageTradeLaneSpeed int
}
