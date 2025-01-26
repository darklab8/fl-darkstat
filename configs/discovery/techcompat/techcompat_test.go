package techcompat

import (
	"testing"

	"github.com/darklab8/fl-darkstat/configs/configs_mapped/parserutils/filefind/file"
	"github.com/darklab8/fl-darkstat/configs/configs_mapped/parserutils/iniload"
	"github.com/darklab8/go-utils/utils/utils_filepath"
	"github.com/darklab8/go-utils/utils/utils_os"
	"github.com/darklab8/go-utils/utils/utils_types"
	"github.com/stretchr/testify/assert"
)

func TestReader(t *testing.T) {
	test_directory := utils_os.GetCurrrentTestFolder()
	fileref := file.NewFile(utils_types.FilePath(utils_filepath.Join(test_directory, "techcompat.cfg")))

	config := Read(iniload.NewLoader(fileref).Scan())

	assert.Greater(t, len(config.Factions), 0)

	for _, faction := range config.Factions {
		faction.ID.Get()
		for _, tech_compat := range faction.TechCompats {
			tech_compat.Nickname.Get()
			tech_compat.Percentage.Get()
		}
	}

	for _, tech_group := range config.TechGroups {
		tech_group.Name.Get()
		for _, item := range tech_group.Items {
			item.Get()
		}
	}

	assert.Equal(t, float64(1), config.GetCompatibilty("drone_miner", "dsy_license_br_m_grp"))
	assert.Equal(t, float64(0.01), config.GetCompatibilty("drone_miner", "dsy_license_vagrants"))
	assert.Equal(t, float64(1), config.GetCompatibilty("nomad_turret01_mark05", "dsy_license_vagrants"))
	assert.Equal(t, float64(0.01), config.GetCompatibilty("nomad_turret01_mark05", ""))
	assert.Equal(t, float64(0.01), config.GetCompatibilty("drone_miner", ""))
	assert.Equal(t, float64(0.9), config.GetCompatibilty("dsy_orderrecon", "dsy_license_srp_07"))

	assert.Equal(t, float64(1), config.GetCompatibilty("item_is_not_listed_in_tech", "someid"))
	assert.Equal(t, float64(0.01), config.GetCompatibilty("compatibility_for_not_equiped_id", ""))

	assert.Equal(t, float64(1.0), config.GetCompatibilty("dsy_no2_cruiser", "dsy_license_nomadguard"))
}
