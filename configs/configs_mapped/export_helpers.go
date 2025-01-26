package configs_mapped

import (
	"fmt"
	"strings"

	"github.com/darklab8/fl-darkstat/configs/configs_mapped/freelancer_mapped/data_mapped/universe_mapped"
	"github.com/darklab8/fl-darkstat/configs/configs_settings"
)

func (configs *MappedConfigs) GetInfocardName(ids_name int, nickname string) string {
	if configs_settings.Env.FallbackInfonamesToNickname {
		return fmt.Sprintf("[%s]", nickname)
	}

	if infoname, ok := configs.Infocards.Infonames[ids_name]; ok {
		return strings.ReplaceAll(string(infoname), "\r", "")
	} else {
		return fmt.Sprintf("[%s]", nickname)
	}
}

func (configs *MappedConfigs) GetRegionName(system *universe_mapped.System) string {
	var Region string
	system_infocard_Id := system.Ids_info.Get()
	if value, ok := configs.Infocards.Infocards[system_infocard_Id]; ok {
		if len(value.Lines) > 0 {
			Region = value.Lines[0]
		}
	}

	if strings.Contains(Region, "Sometimes limbo") && len(Region) > 11 {
		Region = Region[:20] + "..."
	}
	return Region
}
