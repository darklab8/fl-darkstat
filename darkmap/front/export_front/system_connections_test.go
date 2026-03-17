package export_front

import (
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSystemsConnections(t *testing.T) {
	ctx := context.Background()
	export := NewExport(ctx)
	systems := ExportSystems(export.Mapped)
	graph := export.GetSystemConnections(systems)

	for _, system := range graph.Systems {
		if system.Nickname != "hi10" {
			continue
		}
		fmt.Println(system.Nickname, system.Reachable)
	}

	assert.Greater(t, len(systems), 0)
}
