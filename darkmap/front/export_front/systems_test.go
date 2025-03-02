package export_front

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSystemsExport(t *testing.T) {
	export := NewExport()
	systems := exportSystems(export.Mapped)

	assert.Greater(t, len(systems), 0)
}
