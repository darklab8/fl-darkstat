package export_front

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSystemsExport(t *testing.T) {
	ctx := context.Background()
	export := NewExport(ctx)
	systems := exportSystems(export.Mapped)

	assert.Greater(t, len(systems), 0)
}
