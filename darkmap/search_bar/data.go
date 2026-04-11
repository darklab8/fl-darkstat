package search_bar

import (
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
	Tag    string
	Color  string
	Letter string

	SysNick string
	Query   string
}

func NewEntry(
	Name string,
	Tag string,
	Color string,
	Letter string,

	SysNick string,
	Query string,

) Entry {
	return Entry{
		Name:    Name,
		Tag:     Tag,
		Color:   Color,
		Letter:  Letter,
		SysNick: SysNick,
		Query:   Query,
	}
}

func (e Entry) SearchIndex() string {
	return strings.ToLower(e.Name)
}
