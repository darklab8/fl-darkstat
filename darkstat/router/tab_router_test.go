package router

import (
	"testing"

	"github.com/darklab8/fl-darkcore/darkcore/builder"
)

type SomeTabMode int64

func (s SomeTabMode) ToInt() int64 { return int64(s) }

const (
	TabMode1 SomeTabMode = iota
	TabMode2
)

func TestTabRouter(t *testing.T) {
	router := NewTabRouter[string](&builder.Builder{}, []string{}, func(items []string) []string { return items })
	_ = router
}
