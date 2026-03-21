package export_map

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestExporting(t *testing.T) {
	ctx := context.Background()
	export := NewExport(ctx)
	export.Export(ctx)

	assert.Greater(t, len(export.Systems), 0)
}
