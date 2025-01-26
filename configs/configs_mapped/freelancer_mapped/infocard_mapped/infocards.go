package infocard_mapped

import (
	"regexp"
	"strconv"
	"strings"
	"sync"

	"github.com/darklab8/fl-darkstat/configs/configs_mapped/freelancer_mapped/exe_mapped"
	"github.com/darklab8/fl-darkstat/configs/configs_mapped/freelancer_mapped/infocard_mapped/infocard"
	"github.com/darklab8/fl-darkstat/configs/configs_mapped/parserutils/filefind"
	"github.com/darklab8/fl-darkstat/configs/configs_mapped/parserutils/filefind/file"
	"github.com/darklab8/fl-darkstat/configs/configs_mapped/parserutils/inireader"

	"github.com/darklab8/fl-darkstat/configs/configs_settings/logus"
	logus1 "github.com/darklab8/fl-darkstat/configs/configs_settings/logus"
	"github.com/darklab8/go-typelog/typelog"
)

const (
	FILENAME          = "infocards.txt"
	FILENAME_FALLBACK = "infocards.xml"
)

/* For darklint */
func ReadFromTextFile(input_file *file.File) (*infocard.Config, error) {
	frelconfig := infocard.NewConfig()

	lines, err := input_file.ReadLines()
	if logus.Log.CheckError(err, "unable to read file in ReadFromTestFile") {
		return nil, err
	}

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
	return frelconfig, nil
}

var InfocardReader *regexp.Regexp

func init() {
	inireader.InitRegexExpression(&InfocardReader, `^([0-9]+)[= ]+(.*)$`)
}

func ReadFromDiscoServerConfig(input_file *file.File) (*infocard.Config, error) {
	frelconfig := infocard.NewConfig()

	lines, err := input_file.ReadLines()
	if logus.Log.CheckError(err, "unable to read file in ReadFromTestFile") {
		return nil, err
	}

	for _, line := range lines {

		line_match := InfocardReader.FindStringSubmatch(line)
		if len(line_match) <= 0 {
			continue
		}

		id, err := strconv.Atoi(line_match[1])
		if err != nil {
			continue
		}

		content := line_match[2]

		var kind infocard.RecordKind = infocard.TYPE_NAME
		if strings.Contains(content, "<?xml") {
			kind = infocard.TYPE_INFOCAD
		}

		switch kind {
		case infocard.TYPE_NAME:
			frelconfig.Infonames[id] = infocard.Infoname(content)
		case infocard.TYPE_INFOCAD:
			frelconfig.Infocards[id] = infocard.NewInfocard(content)
		default:
			logus1.Log.Fatal(
				"unrecognized object name in infocards.txt",
				typelog.Any("id", id),
				typelog.Any("content", content),
			)
		}
	}
	return frelconfig, nil
}

func Read(filesystem *filefind.Filesystem, freelancer_ini *exe_mapped.Config, input_file *file.File) (*infocard.Config, error) {
	var config *infocard.Config

	config = exe_mapped.GetAllInfocards(filesystem, freelancer_ini.GetDlls())

	// TODO Read from text only for darklint // move this logic to it
	// var err error
	// if input_file != nil {
	// 	config, err = ReadFromTextFile(input_file)
	// 	if logus.Log.CheckError(err, "unable to read infocards", typelog.OptError(err)) {
	// 		return config, err
	// 	}
	// } else {

	// }

	if input_file != nil {
		config_override, err := ReadFromDiscoServerConfig(input_file)
		logus.Log.CheckPanic(err, "unable to read infocards", typelog.OptError(err))

		for key, value := range config_override.Infocards {
			config.Infocards[key] = value
		}
		for key, value := range config_override.Infonames {
			config.Infonames[key] = value
		}
	}

	var wg sync.WaitGroup
	for _, card := range config.Infocards {
		wg.Add(1)
		go func(parsed_infocard *infocard.Infocard) {
			card.Lines, _ = card.XmlToText()
			wg.Done()
		}(card)
	}
	wg.Wait()
	return config, nil
}
