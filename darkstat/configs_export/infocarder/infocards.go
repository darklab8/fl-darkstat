/* Infocarder is thread safe version of infocards management for configs_export */
package infocarder

import (
	"strings"
	"sync"
)

type InfocardKey string

type Infocarder struct {
	Mutex     sync.RWMutex
	Infocards Infocards
}

func NewInfocarder() Infocarder {
	return Infocarder{
		Infocards: make(Infocards),
	}
}

type InfocardPhrase struct {
	Phrase string  `json:"phrase"  validate:"required"`
	Link   *string `json:"link"`
	Bold   bool    `json:"bold"  validate:"required"`
}

type InfocardLine struct {
	Phrases []InfocardPhrase `json:"phrases"  validate:"required"`
}

func (i InfocardLine) ToStr() string {
	var sb strings.Builder
	for _, phrase := range i.Phrases {
		sb.WriteString(phrase.Phrase)
	}
	return sb.String()
}

func NewInfocardSimpleLine(line string) InfocardLine {
	return InfocardLine{Phrases: []InfocardPhrase{{Phrase: line}}}
}

func NewInfocardBuilder() InfocardBuilder {
	return InfocardBuilder{}
}
func (i *InfocardBuilder) WriteLine(phrases ...InfocardPhrase) {
	i.Lines = append(i.Lines, InfocardLine{Phrases: phrases})
}
func (i *InfocardBuilder) WriteLineStr(phrase_strs ...string) {
	var phrases []InfocardPhrase
	for _, phrase := range phrase_strs {
		phrases = append(phrases, InfocardPhrase{Phrase: phrase})
	}
	i.Lines = append(i.Lines, InfocardLine{Phrases: phrases})
}

type InfocardBuilder struct {
	Lines Infocard
}

type Infocard []InfocardLine

func (i Infocard) StringsJoin(delimiter string) string {
	var sb strings.Builder

	for _, line := range i {
		for _, phrase := range line.Phrases {
			sb.WriteString(phrase.Phrase)
		}
		sb.WriteString(delimiter)
	}
	return sb.String()
}

type Infocards map[InfocardKey]Infocard
