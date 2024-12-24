package types

import (
	"context"
	"time"

	"github.com/darklab8/fl-configs/configs/cfgtype"
	"github.com/darklab8/fl-configs/configs/configs_export"
	"github.com/darklab8/fl-configs/configs/discovery/techcompat"
	"github.com/darklab8/fl-darkcore/darkcore/core_types"
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
	SiteRoot       string
	StaticRoot     string
	Heading        string
	Timestamp      time.Time
	TractorTabName string
}

func (g *GlobalParams) GetBuildPath() utils_types.FilePath {
	return g.Buildpath
}

func (g *GlobalParams) GetStaticRoot() string {
	return g.StaticRoot
}

var check core_types.GlobalParamsI = &GlobalParams{}

func GetCtx(ctx context.Context) *GlobalParams {
	return ctx.Value(core_types.GlobalParamsCtxKey).(*GlobalParams)
}

type FLSRData struct {
}

type DiscoveryData struct {
	ShowDisco    bool
	Ids          []configs_export.Tractor
	TractorsByID map[cfgtype.TractorID]configs_export.Tractor
	Config       *techcompat.Config
	LatestPatch  autopatcher.Patch

	Infocards configs_export.Infocards

	OrderedTechcompat configs_export.TechCompatOrderer
}

type SharedData struct {
	DiscoveryData
	FLSRData
}
