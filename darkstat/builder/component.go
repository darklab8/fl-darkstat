package builder

import (
	"bytes"
	"context"

	"github.com/darklab8/fl-darkstat/darkstat/common/types"
	"github.com/yosssi/gohtml"

	"github.com/a-h/templ"
	"github.com/darklab8/go-utils/goutils/utils/utils_filepath"
	"github.com/darklab8/go-utils/goutils/utils/utils_types"
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

func (h *Component) Write(gp types.GlobalParams, filesystem *Filesystem) {
	buf := bytes.NewBuffer([]byte{})

	gp.Pagepath = string(h.pagepath)

	realpath := utils_filepath.Join(gp.Buildpath, h.pagepath)

	h.templ_comp.Render(context.WithValue(context.Background(), types.GlobalParamsCtxKey, gp), buf)

	// Usage of gohtml is not obligatory, but nice touch simplifying debugging view.
	filesystem.WriteToMem(realpath, gohtml.FormatBytes(buf.Bytes()))
}
