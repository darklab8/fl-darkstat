package builder

import (
	"bytes"
	"context"
	"fmt"
	"time"

	"github.com/darklab8/fl-darkstat/darkcore/core_types"
	"github.com/darklab8/fl-darkstat/darkcore/settings/logus"
	"github.com/darklab8/fl-darkstat/darkcore/settings/traces"
	"github.com/darklab8/go-utils/typelog"

	"github.com/a-h/templ"
	"github.com/darklab8/go-utils/utils/utils_filepath"
	"github.com/darklab8/go-utils/utils/utils_types"
)

var (
	Log *typelog.Logger = logus.Log.WithScope("darkcore.builder.component")
)

type Component struct {
	pagepath   utils_types.FilePath
	templ_comp templ.Component
}

func NewComponent(
	pagepath utils_types.FilePath,
	templ_comp templ.Component,
) *Component {
	return &Component{
		pagepath:   pagepath,
		templ_comp: templ_comp,
	}
}

type WriteResult struct {
	realpath utils_types.FilePath
	bytes    []byte
}

func (h *Component) GetPagePath(gp Params) utils_types.FilePath {
	return utils_filepath.Join(gp.GetBuildPath(), h.pagepath)
}

func (h *Component) Write(ctx context.Context, gp Params) WriteResult {
	ctx, span := traces.Tracer.Start(ctx, "component-write")
	defer span.End()

	buf := bytes.NewBuffer([]byte{})
	buf.Write([]byte(fmt.Sprintf("<!--ts %s-->\n", time.Now().Format("2006-01-02T15:04:05.999Z"))))
	// gp.Pagepath = string(h.pagepath)

	err := h.templ_comp.Render(context.WithValue(ctx, core_types.GlobalParamsCtxKey, gp), buf)
	Log.CheckPanic(err, "failed to write component")

	// Usage of gohtml is not obligatory, but nice touch simplifying debugging view.
	return WriteResult{
		realpath: h.GetPagePath(gp),
		bytes:    buf.Bytes(),
	}
}
