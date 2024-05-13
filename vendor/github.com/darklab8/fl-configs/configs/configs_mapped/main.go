/*
Tool to parse freelancer configs
*/
package configs_mapped

import (
	"sync"

	"github.com/darklab8/fl-configs/configs/configs_mapped/freelancer_mapped/data_mapped/const_mapped"
	"github.com/darklab8/fl-configs/configs/configs_mapped/freelancer_mapped/data_mapped/equipment_mapped"
	"github.com/darklab8/fl-configs/configs/configs_mapped/freelancer_mapped/data_mapped/initialworld"
	"github.com/darklab8/fl-configs/configs/configs_mapped/freelancer_mapped/data_mapped/interface_mapped"
	"github.com/darklab8/fl-configs/configs/configs_mapped/freelancer_mapped/data_mapped/missions_mapped/empathy_mapped"
	"github.com/darklab8/fl-configs/configs/configs_mapped/freelancer_mapped/data_mapped/missions_mapped/mbases_mapped"
	"github.com/darklab8/fl-configs/configs/configs_mapped/freelancer_mapped/data_mapped/ship_mapped"
	"github.com/darklab8/fl-configs/configs/configs_mapped/freelancer_mapped/data_mapped/universe_mapped"
	"github.com/darklab8/fl-configs/configs/configs_mapped/freelancer_mapped/data_mapped/universe_mapped/systems_mapped"
	"github.com/darklab8/fl-configs/configs/configs_mapped/freelancer_mapped/exe_mapped"
	"github.com/darklab8/fl-configs/configs/configs_mapped/freelancer_mapped/infocard_mapped"
	"github.com/darklab8/fl-configs/configs/configs_mapped/freelancer_mapped/infocard_mapped/infocard"
	"github.com/darklab8/fl-configs/configs/configs_mapped/parserutils/filefind"
	"github.com/darklab8/fl-configs/configs/configs_mapped/parserutils/filefind/file"
	"github.com/darklab8/fl-configs/configs/configs_mapped/parserutils/iniload"
	"github.com/darklab8/fl-configs/configs/configs_mapped/parserutils/semantic"
	"github.com/darklab8/fl-configs/configs/settings/logus"

	"github.com/darklab8/fl-configs/configs/configs_mapped/freelancer_mapped/data_mapped/equipment_mapped/equip_mapped"
	"github.com/darklab8/fl-configs/configs/configs_mapped/freelancer_mapped/data_mapped/equipment_mapped/market_mapped"
	"github.com/darklab8/fl-configs/configs/configs_mapped/freelancer_mapped/data_mapped/equipment_mapped/weaponmoddb"
	"github.com/darklab8/fl-configs/configs/discovery/discoprices"
	"github.com/darklab8/fl-configs/configs/discovery/techcompat"

	"github.com/darklab8/go-utils/goutils/utils"
	"github.com/darklab8/go-utils/goutils/utils/time_measure"
	"github.com/darklab8/go-utils/goutils/utils/utils_logus"
	"github.com/darklab8/go-utils/goutils/utils/utils_types"
)

type DiscoveryConfig struct {
	Techcompat *techcompat.Config
	Prices     *discoprices.Config
}

type MappedConfigs struct {
	FreelancerINI *exe_mapped.Config

	Universe_config *universe_mapped.Config
	Systems         *systems_mapped.Config

	Market   *market_mapped.Config
	Equip    *equip_mapped.Config
	Goods    *equipment_mapped.Config
	Shiparch *ship_mapped.Config

	InfocardmapINI *interface_mapped.Config
	Infocards      *infocard.Config
	InitialWorld   *initialworld.Config
	Empathy        *empathy_mapped.Config
	MBases         *mbases_mapped.Config
	Consts         *const_mapped.Config
	WeaponMods     *weaponmoddb.Config

	Discovery *DiscoveryConfig
}

func NewMappedConfigs() *MappedConfigs {
	return &MappedConfigs{}
}

func getConfigs(filesystem *filefind.Filesystem, paths []*semantic.Path) []*iniload.IniLoader {
	return utils.CompL(paths, func(x *semantic.Path) *iniload.IniLoader {
		return iniload.NewLoader(filesystem.GetFile(utils_types.FilePath(x.FileName())))
	})
}

func (p *MappedConfigs) Read(file1path utils_types.FilePath) *MappedConfigs {
	logus.Log.Info("Parse START for FreelancerFolderLocation=", utils_logus.FilePath(file1path))
	filesystem := filefind.FindConfigs(file1path)
	p.FreelancerINI = exe_mapped.Read(iniload.NewLoader(filesystem.GetFile(exe_mapped.FILENAME_FL_INI)).Scan())

	files_goods := getConfigs(filesystem, p.FreelancerINI.Goods)
	files_market := getConfigs(filesystem, p.FreelancerINI.Markets)
	files_equip := getConfigs(filesystem, p.FreelancerINI.Equips)
	files_shiparch := getConfigs(filesystem, p.FreelancerINI.Ships)
	file_universe := iniload.NewLoader(filesystem.GetFile(universe_mapped.FILENAME))
	file_interface := iniload.NewLoader(filesystem.GetFile(interface_mapped.FILENAME_FL_INI))
	file_initialworld := iniload.NewLoader(filesystem.GetFile(initialworld.FILENAME))
	file_empathy := iniload.NewLoader(filesystem.GetFile(empathy_mapped.FILENAME))
	file_mbases := iniload.NewLoader(filesystem.GetFile(mbases_mapped.FILENAME))
	file_consts := iniload.NewLoader(filesystem.GetFile(const_mapped.FILENAME))
	file_weaponmoddb := iniload.NewLoader(filesystem.GetFile(weaponmoddb.FILENAME))

	all_files := append(files_goods, files_market...)
	all_files = append(all_files, files_equip...)
	all_files = append(all_files, files_shiparch...)
	all_files = append(all_files,
		file_universe,
		file_interface,
		file_initialworld,
		file_empathy,
		file_mbases,
		file_consts,
		file_weaponmoddb,
	)

	var file_techcompat *iniload.IniLoader
	var file_prices *iniload.IniLoader
	if techcom := filesystem.GetFile("launcherconfig.xml"); techcom != nil {
		p.Discovery = &DiscoveryConfig{}
		file_techcompat = iniload.NewLoader(file.NewWebFile("https://discoverygc.com/gameconfigpublic/techcompat.cfg"))
		file_prices = iniload.NewLoader(file.NewWebFile("https://discoverygc.com/gameconfigpublic/prices.cfg"))
		all_files = append(all_files, file_techcompat, file_prices)
	}

	time_measure.TimeMeasure(func(m *time_measure.TimeMeasurer) {
		var wg sync.WaitGroup
		for _, file := range all_files {
			wg.Add(1)
			go func(file *iniload.IniLoader) {
				file.Scan()
				wg.Done()
			}(file)
		}
		wg.Wait()
	}, time_measure.WithMsg("Scanned ini loaders"))

	time_measure.TimeMeasure(func(m *time_measure.TimeMeasurer) {
		p.Universe_config = universe_mapped.Read(file_universe)

		p.Systems = systems_mapped.Read(p.Universe_config, filesystem)

		p.Market = market_mapped.Read(files_market)
		p.Equip = equip_mapped.Read(files_equip)
		p.Goods = equipment_mapped.Read(files_goods)
		p.Shiparch = ship_mapped.Read(files_shiparch)

		p.InfocardmapINI = interface_mapped.Read(file_interface)

		var infocards_override *file.File
		if p.Discovery != nil {
			infocards_override = filesystem.GetFile("temp.disco.infocards.txt")

			if infocards_override == nil {
				infocards_override = iniload.NewLoader(file.NewWebFile("https://discoverygc.com/gameconfigpublic/infocard_overrides.cfg"))
			}
		}

		p.Infocards, _ = infocard_mapped.Read(filesystem, p.FreelancerINI, infocards_override)

		p.InitialWorld = initialworld.Read(file_initialworld)
		p.Empathy = empathy_mapped.Read(file_empathy)
		p.MBases = mbases_mapped.Read(file_mbases)
		p.Consts = const_mapped.Read(file_consts)
		p.WeaponMods = weaponmoddb.Read(file_weaponmoddb)

		if p.Discovery != nil {
			p.Discovery.Techcompat = techcompat.Read(file_techcompat)
			p.Discovery.Prices = discoprices.Read(file_prices)
		}
	}, time_measure.WithMsg("Mapped stuff"))

	logus.Log.Info("Parse OK for FreelancerFolderLocation=", utils_logus.FilePath(file1path))

	return p
}

type IsDruRun bool

func (p *MappedConfigs) Write(is_dry_run IsDruRun) {
	// write
	files := []*file.File{}

	files = append(files, p.Universe_config.Write())
	files = append(files, p.Systems.Write()...)
	files = append(files, p.Market.Write()...)
	files = append(files, p.Equip.Write()...)
	files = append(files, p.Goods.Write()...)
	files = append(files, p.Shiparch.Write()...)
	files = append(files, p.InfocardmapINI.Write())
	files = append(files, p.InitialWorld.Write())
	files = append(files, p.Empathy.Write())
	files = append(files, p.MBases.Write())
	files = append(files, p.Consts.Write())
	files = append(files, p.WeaponMods.Write())

	if is_dry_run {
		return
	}

	for _, file := range files {
		file.WriteLines()
	}
}
