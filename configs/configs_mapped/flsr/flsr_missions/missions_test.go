package flsr_missions

import (
	"fmt"
	"testing"

	"github.com/darklab8/fl-darkstat/configs/configs_mapped/parserutils/filefind"
	"github.com/darklab8/fl-darkstat/configs/configs_mapped/parserutils/iniload"
	"github.com/darklab8/fl-darkstat/configs/configs_settings/logus"
	"github.com/darklab8/fl-darkstat/configs/tests"
	"github.com/darklab8/fl-darkstat/darkstat/settings"
	"github.com/stretchr/testify/assert"
)

func TestReader(t *testing.T) {
	fileref := tests.FixtureFileFind()
	if fileref.GetFile("flsr-launcher.ini") == nil {
		return
	}

	file1path := settings.Env.FreelancerFolder
	filesystem := filefind.FindConfigs(file1path.Join("EXE", "flhook_plugins", "missions"))
	if len(filesystem.Files) == 0 {
		logus.Log.Panic("expected to find FLSR/flhoook_plugins/missions, but found none :|")
	}

	loaded_files := []*iniload.IniLoader{}
	for _, file := range filesystem.Files {
		loaded_files = append(loaded_files, iniload.NewLoader(file).Scan())
	}

	config := Read(loaded_files)

	assert.Greater(t, len(config.Missions), 0, "expected finding missions")

	for _, mission := range config.Missions {
		active_missions, _ := mission.InitState.GetValue()
		if len(mission.Solars) > 0 {
			if active_missions {
				fmt.Println("mission.Solars=", len(mission.Solars), " initstate=", active_missions, " nickname=", mission.Nickname.Get())
			}
		}
	}

}
