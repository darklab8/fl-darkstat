package search_bar

import (
	"fmt"
	"strings"

	_ "embed"

	"github.com/darklab8/fl-darkstat/darkcore/core_types"
)

//go:embed search_bar.css
var SearchBarContent string

var SearchBarCss = core_types.StaticFile{
	Content:  SearchBarContent,
	Filename: "search_bar.css",
	Kind:     core_types.StaticFileCSS,
}

type Entry struct {
	Name   string
	Kind   string
	Tag    string
	Color  string
	Letter string

	SysNick string
	Query   string
}

func NewEntry(
	Name string,
	Kind string,
	Tag string,
	Color string,
	Letter string,

	SysNick string,
	Query string,

) Entry {
	return Entry{
		Name:    Name,
		Kind:    Kind,
		Tag:     Tag,
		Color:   Color,
		Letter:  Letter,
		SysNick: SysNick,
		Query:   Query,
	}
}

func (e Entry) SearchIndex() string {
	return strings.ToLower(fmt.Sprintf("%s: %s", e.Kind, e.Name))
}
