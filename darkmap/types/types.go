package types

import (
	"context"
	"time"

	"github.com/darklab8/fl-darkstat/configs/discovery/techcompat"
	"github.com/darklab8/fl-darkstat/darkcore/core_types"
	"github.com/darklab8/fl-darkstat/darkstat/configs_export"
	"github.com/darklab8/fl-data-discovery/autopatcher"
	"github.com/darklab8/go-utils/utils/utils_types"
)

type GlobalParams struct {
	Buildpath         utils_types.FilePath
	SiteRoot          string
	StaticRoot        string
	OppositeThemeRoot string
	Timestamp         time.Time
}

func (g GlobalParams) GetBuildPath() utils_types.FilePath {
	return g.Buildpath
}

func (g GlobalParams) GetStaticRoot() string {
	return g.StaticRoot
}

var check core_types.GlobalParamsI = GlobalParams{}

func GetCtx(ctx context.Context) GlobalParams {
	return ctx.Value(core_types.GlobalParamsCtxKey).(GlobalParams)
}

type DiscoveryIDs struct {
	Show        bool
	Ids         []configs_export.Tractor
	Config      *techcompat.Config
	LatestPatch autopatcher.Patch
}
