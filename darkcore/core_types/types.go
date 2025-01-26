package core_types

import (
	"context"
)

type CtxKey string

const GlobalParamsCtxKey CtxKey = "global_params"

type Url string

type StaticFileKind int64

const (
	StaticFileUnknown StaticFileKind = iota // default unkonwn
	StaticFileJS
	StaticFileIco
	StaticFileCSS
)

type StaticFile struct {
	Content  string
	Filename string
	Kind     StaticFileKind
}

type GlobalParamsI interface {
	GetStaticRoot() string
}

func GetCtx(ctx context.Context) GlobalParamsI {
	return ctx.Value(GlobalParamsCtxKey).(GlobalParamsI)
}
