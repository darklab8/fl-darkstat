package inireader

import (
	"testing"

	"github.com/darklab8/fl-darkstat/configs/configs_mapped/parserutils/filefind"
	"github.com/darklab8/fl-darkstat/configs/configs_mapped/parserutils/filefind/file"
	"github.com/darklab8/fl-darkstat/configs/tests"
	"github.com/darklab8/go-utils/utils/utils_os"
	"github.com/stretchr/testify/assert"
)

func TestReader(t *testing.T) {
	fileref := tests.FixtureFileFind().GetFile("market_ships.ini")
	config := Read(fileref)

	assert.Greater(t, len(config.Sections), 0, "market ships sections were not scanned")
}

func TestReaderWithBOMFails(t *testing.T) {

	defer func() {
		InitRegexExpression(&regexSection, regexSectionRegExp)
	}()
	InitRegexExpression(&regexSection, `^\[.*\]`)

	fs := filefind.FindConfigs(utils_os.GetCurrrentTestFolder())
	fileref := fs.GetFile("li05_with_bom.ini")

	var crashed bool = false
	func() {
		defer func() {
			if r := recover(); r != nil {
				crashed = true
			}
		}()
		Read(fileref)
	}()

	assert.True(t, crashed, "with BOM we Crash.")
}

func TestReaderWithBOMPasses(t *testing.T) {

	fs := filefind.FindConfigs(utils_os.GetCurrrentTestFolder())
	fileref := fs.GetFile("li05_with_bom.ini")
	config := Read(fileref)

	assert.Greater(t, len(config.Sections), 0, "market ships sections were not scanned")
}

func TestReadScientificNotation(t *testing.T) {

	fs := filefind.FindConfigs(utils_os.GetCurrrentTestFolder())
	fileref := fs.GetFile("hud.ini")
	config := Read(fileref)

	assert.Greater(t, len(config.Sections), 0, "expected not zero section")
}

func TestCommentsHandling(t *testing.T) {

	fs := filefind.FindConfigs(utils_os.GetCurrrentTestFolder())
	fileref := fs.GetFile("comments.ini")
	config := Read(fileref)

	assert.Greater(t, len(config.Sections), 0, "expected not zero section")

}

func TestCommentsHandling2(t *testing.T) {

	fs := filefind.FindConfigs(utils_os.GetCurrrentTestFolder())
	fileref := fs.GetFile("comments2.ini")
	config := Read(fileref)

	assert.Greater(t, len(config.Sections), 0, "expected not zero section")

}

func TestCommentsHandling3(t *testing.T) {

	fs := filefind.FindConfigs(utils_os.GetCurrrentTestFolder())
	fileref := fs.GetFile("comments3.ini")
	config := Read(fileref)

	assert.Greater(t, len(config.Sections), 0, "expected not zero section")

}

func TestCommentsHandling4(t *testing.T) {

	fs := filefind.FindConfigs(utils_os.GetCurrrentTestFolder())
	fileref := fs.GetFile("comments4.ini")
	config := Read(fileref)

	assert.Greater(t, len(config.Sections), 0, "expected not zero section")

}

func TestCommentsHandling5(t *testing.T) {

	fs := filefind.FindConfigs(utils_os.GetCurrrentTestFolder())
	fileref := fs.GetFile("comments5.ini")
	config := Read(fileref)

	write_file := file.NewFile(utils_os.GetCurrrentTestFolder().Join("comments5_rendered.ini"))
	config.Write(write_file)
	write_file.WriteLines()
	assert.Greater(t, len(config.Sections), 0, "expected not zero section")
}

func TestCommentsHandling6(t *testing.T) {

	fs := filefind.FindConfigs(utils_os.GetCurrrentTestFolder())
	fileref := fs.GetFile("comments6.ini")
	config := Read(fileref)

	write_file := file.NewFile(utils_os.GetCurrrentTestFolder().Join("comments6_rendered.ini"))
	config.Write(write_file)
	write_file.WriteLines()
	assert.Greater(t, len(config.Sections), 0, "expected not zero section")
}
