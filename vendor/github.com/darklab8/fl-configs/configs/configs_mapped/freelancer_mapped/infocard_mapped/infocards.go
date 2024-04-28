package infocard_mapped

import (
	"strconv"
	"sync"

	"github.com/darklab8/fl-configs/configs/configs_mapped/freelancer_mapped/exe_mapped"
	"github.com/darklab8/fl-configs/configs/configs_mapped/freelancer_mapped/infocard_mapped/infocard"
	"github.com/darklab8/fl-configs/configs/configs_mapped/parserutils/filefind"
	"github.com/darklab8/fl-configs/configs/configs_mapped/parserutils/filefind/file"
	logus1 "github.com/darklab8/fl-configs/configs/settings/logus"
	"github.com/darklab8/go-typelog/typelog"
)

const (
	FILENAME          = "infocards.txt"
	FILENAME_FALLBACK = "infocards.xml"
)

func ReadFromTextFile(input_file *file.File) *infocard.Config {
	frelconfig := infocard.NewConfig()

	lines := input_file.ReadLines()

	for index := 0; index < len(lines); index++ {

		id, err := strconv.Atoi(lines[index])
		if err != nil {
			continue
		}
		name := lines[index+1]
		content := lines[index+2]
		index += 3

		switch infocard.RecordKind(name) {
		case infocard.TYPE_NAME:
			frelconfig.Infonames[id] = infocard.Infoname(content)
		case infocard.TYPE_INFOCAD:
			frelconfig.Infocards[id] = infocard.NewInfocard(content)
		default:
			logus1.Log.Fatal(
				"unrecognized object name in infocards.txt",
				typelog.Any("id", id),
				typelog.Any("name", name),
				typelog.Any("content", content),
			)
		}
	}
	return frelconfig
}

func Read(filesystem *filefind.Filesystem, freelancer_ini *exe_mapped.Config, input_file *file.File) *infocard.Config {
	var config *infocard.Config
	if input_file != nil {
		config = ReadFromTextFile(input_file)
	} else {
		config = exe_mapped.GetAllInfocards(filesystem, freelancer_ini.GetDlls())
	}

	var wg sync.WaitGroup
	for _, card := range config.Infocards {
		wg.Add(1)
		go func(card *infocard.Infocard) {
			card.Lines, _ = card.XmlToText()
			wg.Done()
		}(card)
	}
	wg.Wait()
	return config
}
