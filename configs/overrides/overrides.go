package overrides

import (
	"os"

	"github.com/darklab8/fl-darkstat/configs/configs_settings/logus"
	"github.com/darklab8/go-utils/utils/utils_types"
	"gopkg.in/yaml.v3"
)

const FILENAME = "overrides.fl_configs.yml"

type Overrides struct {
	SystemTravelSpeedMultipliers map[string]float64 `yaml:"system_travel_speed_multilpliers"`
}

type InfocardRegion string

func (o Overrides) GetSystemSpeedMultiplier(system_nickname string) float64 {
	if value, ok := o.SystemTravelSpeedMultipliers[system_nickname]; ok {
		return value
	} else {
		return 1.0
	}
}

func Read(filepath utils_types.FilePath) Overrides {
	var config Overrides
	config.SystemTravelSpeedMultipliers = make(map[string]float64)

	data, err := os.ReadFile(filepath.ToString())
	logus.Log.CheckWarn(err, "overrides for fl configs is not found")

	err = yaml.Unmarshal(data, &config)

	logus.Log.CheckPanic(err, "failed inmarshaling yaml file for overrides")

	return config
}
