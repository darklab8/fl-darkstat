package export_map

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSystemsExport(t *testing.T) {
	ctx := context.Background()
	export := NewExport(ctx)
	systems := export.ExportSystems(export.Mapped)

	assert.Greater(t, len(systems), 0)
}
