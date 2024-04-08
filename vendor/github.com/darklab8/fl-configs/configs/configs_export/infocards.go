package configs_export

import (
	"github.com/darklab8/fl-configs/configs/configs_mapped/freelancer_mapped/infocard_mapped/infocard"
	"github.com/darklab8/fl-configs/configs/settings/logus"
	"github.com/darklab8/go-typelog/typelog"
)

type Infocard struct {
	Lines []string
}

type InfocardKey string

// This abstraction optimizes performance
// for infocards extracting, plus adds reusage
type InfocardsParser struct {
	infocards *infocard.Config
	data      map[InfocardKey]*Infocard

	results       chan *InfocardResult
	awaited_count int
	recalculated  int
}

func NewInfocardsParser(Infocards *infocard.Config) *InfocardsParser {
	return &InfocardsParser{
		infocards: Infocards,
		results:   make(chan *InfocardResult),
		data:      make(map[InfocardKey]*Infocard),
	}
}

type InfocardResult struct {
	Key InfocardKey
	Infocard
}

func (i *InfocardsParser) Set(key InfocardKey, ids ...int) {
	_, exists := i.data[key]
	if exists {
		logus.Log.Debug("such infocard is already parsed", typelog.Any("key", key), typelog.Any("ids", ids))
		return
	}

	i.awaited_count++

	go func() {
		result := &InfocardResult{Key: key}
		for _, id := range ids {
			infocard, infocard_exists := i.infocards.Infocards[id]
			if !infocard_exists {
				logus.Log.Debug("infocard with such id does not exist", typelog.Int("id", id), typelog.Any("key", key), typelog.Any("ids", ids))
				continue
			}
			infocard_lines, err := infocard.XmlToText()
			logus.Log.CheckError(err, "failed to xml infocard")
			result.Lines = append(result.Lines, infocard_lines...)

		}
		i.results <- result
	}()
}

func (s *InfocardsParser) makeReady() {
	for i := 0; i < s.awaited_count; i++ {
		result := <-s.results

		_, exists := s.data[result.Key]
		if exists {
			s.recalculated++
			logus.Log.Debug("you recalculated infocard!", typelog.Any("key", result.Key), typelog.Int("recalculated", s.recalculated))

		}

		s.data[result.Key] = &result.Infocard
	}
	s.awaited_count = 0
}

func (s *InfocardsParser) Get() map[InfocardKey]*Infocard {
	s.makeReady()
	close(s.results)
	for result := range s.results {
		s.data[result.Key] = &result.Infocard
		logus.Log.Warn("you received results after channel closing!")
	}
	if s.recalculated > 0 {
		logus.Log.Warn("you recalculated infocards", typelog.Int("count", s.recalculated))
	}
	return s.data
}
