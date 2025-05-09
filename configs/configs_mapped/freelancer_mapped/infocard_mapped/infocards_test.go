package infocard_mapped

import (
	"fmt"
	"strings"
	"testing"

	"github.com/darklab8/fl-darkstat/configs/configs_mapped/freelancer_mapped/exe_mapped"
	"github.com/darklab8/fl-darkstat/configs/configs_mapped/parserutils/filefind"
	"github.com/darklab8/fl-darkstat/configs/tests"
	"github.com/darklab8/go-utils/utils/utils_os"
	"github.com/stretchr/testify/assert"
)

// Not used any longer?
func TestReader(t *testing.T) {
	one_file_filesystem := filefind.FindConfigs(utils_os.GetCurrrentTestFolder())

	filesystem := tests.FixtureFileFind()

	freelancer_ini := exe_mapped.FixtureFLINIConfig()
	_config, _ := Read(filesystem, freelancer_ini, one_file_filesystem.GetFile("temp.disco.infocards.txt"))
	config := _config.GetUnsafe()

	assert.Greater(t, len(config.Infocards), 0)

	fmt.Println(len(config.Infocards))
	for index, infocard := range config.Infocards {
		if strings.Contains(infocard.GetContent(), "TAU BORDER WORLDS") {
			lines, _ := infocard.XmlToText()
			fmt.Println("index=", index, " ", lines)
		}
	}
}
