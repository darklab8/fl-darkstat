package configs

import (
	"fmt"
	"strings"

	"github.com/darklab8/fl-darkstat/configs/configs_mapped"
	"github.com/darklab8/fl-darkstat/configs/configs_mapped/freelancer_mapped/data_mapped/equipment_mapped/equip_mapped"
	"github.com/darklab8/fl-darkstat/configs/configs_mapped/parserutils/filefind/file"
	"github.com/darklab8/fl-darkstat/configs/configs_mapped/parserutils/iniload"
	"github.com/darklab8/fl-darkstat/configs/configs_settings"
	"github.com/darklab8/fl-darkstat/configs/configs_settings/logus"
	"github.com/darklab8/go-utils/utils/utils_logus"
)

// ExampleImportIniSection demonstrates how to add section to specific section
func Example_importIniSection() {
	// can be having imperfections related to how to handle comments. To improve some day

	freelancer_folder := configs_settings.Env.FreelancerFolder
	mapped := configs_mapped.NewMappedConfigs()
	logus.Log.Debug("scanning freelancer folder", utils_logus.FilePath(freelancer_folder))

	// Reading to ini universal custom format and mapping to ORM objects
	// which have both reading and writing back capabilities
	mapped.Read(freelancer_folder)

	// TODO Adding Section manually
	// var new_section *inireader.Section = &inireader.Section{}
	// mapped_gun := &equip_mapped.Gun{}
	// mapped_gun.Map(new_section)

	// mapped_gun.Nickname.Set("my_gun_nickname")
	// mapped_gun.IdsName.Set(3453453)
	// mapped_gun.HPGunType.Set("some_hp_type")

	// configs.Equip.Files[0].Sections = append(configs.Equip.Files[0].Sections, new_section)

	///////////// Alternatively feeding in memory section entirely

	content := `
[Gun]
nickname = some_gun
hp_gun_type = some_hp_type
`
	memory_file := file.NewMemoryFile(strings.Split(content, "\n"))
	scanned_file := iniload.NewLoader(memory_file).Scan()

	mapped_equip := equip_mapped.Read([]*iniload.IniLoader{scanned_file})
	fmt.Println(mapped_equip.Guns[0].Nickname.Get())
	fmt.Println(mapped_equip.Guns[0].HPGunType.Get())

	// Adding to main Freelancer instance..
	mapped.Equip().Files[0].Sections = append(mapped.Equip().Files[0].Sections, mapped_equip.Guns[0].Model.RenderModel())

	// Write without Dry Run for writing to files modified values back!
	mapped.Write(configs_mapped.IsDruRun(true))
}
