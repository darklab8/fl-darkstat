package types

import (
	"context"
	"time"

	"github.com/darklab8/go-utils/utils/utils_types"
)

type Url string

type Theme int64

type CtxKey string

const (
	ThemeDark Theme = iota
	ThemeLight
)

const GlobalParamsCtxKey CtxKey = "global_params"

type GlobalParams struct {
	Buildpath         utils_types.FilePath
	Theme             Theme
	SiteRoot          string
	StaticRoot        string
	OppositeThemeRoot string
	Heading           string
	Timestamp         time.Time
	TractorTabName    string
}

func (g GlobalParams) GetBuildPath() utils_types.FilePath {
	return g.Buildpath
}

func GetCtx(ctx context.Context) GlobalParams {
	return ctx.Value(GlobalParamsCtxKey).(GlobalParams)
}
