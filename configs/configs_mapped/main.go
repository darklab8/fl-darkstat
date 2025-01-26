/*
Tool to parse freelancer configs
*/
package configs_mapped

import (
	"encoding/json"
	"fmt"
	"os"
	"sync"

	"github.com/darklab8/fl-darkstat/configs/configs_mapped/flsr/flsr_recipes"
	"github.com/darklab8/fl-darkstat/configs/configs_mapped/freelancer_mapped/data_mapped/const_mapped"
	"github.com/darklab8/fl-darkstat/configs/configs_mapped/freelancer_mapped/data_mapped/equipment_mapped"
	"github.com/darklab8/fl-darkstat/configs/configs_mapped/freelancer_mapped/data_mapped/initialworld"
	"github.com/darklab8/fl-darkstat/configs/configs_mapped/freelancer_mapped/data_mapped/interface_mapped"
	"github.com/darklab8/fl-darkstat/configs/configs_mapped/freelancer_mapped/data_mapped/missions_mapped/empathy_mapped"
	"github.com/darklab8/fl-darkstat/configs/configs_mapped/freelancer_mapped/data_mapped/missions_mapped/faction_props_mapped"
	"github.com/darklab8/fl-darkstat/configs/configs_mapped/freelancer_mapped/data_mapped/missions_mapped/mbases_mapped"
	"github.com/darklab8/fl-darkstat/configs/configs_mapped/freelancer_mapped/data_mapped/missions_mapped/npc_ships"
	"github.com/darklab8/fl-darkstat/configs/configs_mapped/freelancer_mapped/data_mapped/rnd_msns_mapped/diff2money"
	"github.com/darklab8/fl-darkstat/configs/configs_mapped/freelancer_mapped/data_mapped/rnd_msns_mapped/npcranktodiff"
	"github.com/darklab8/fl-darkstat/configs/configs_mapped/freelancer_mapped/data_mapped/ship_mapped"
	"github.com/darklab8/fl-darkstat/configs/configs_mapped/freelancer_mapped/data_mapped/solar_mapped/loadouts_mapped"
	"github.com/darklab8/fl-darkstat/configs/configs_mapped/freelancer_mapped/data_mapped/solar_mapped/solararch_mapped"
	"github.com/darklab8/fl-darkstat/configs/configs_mapped/freelancer_mapped/data_mapped/universe_mapped"
	"github.com/darklab8/fl-darkstat/configs/configs_mapped/freelancer_mapped/data_mapped/universe_mapped/systems_mapped"
	"github.com/darklab8/fl-darkstat/configs/configs_mapped/freelancer_mapped/exe_mapped"
	"github.com/darklab8/fl-darkstat/configs/configs_mapped/freelancer_mapped/infocard_mapped"
	"github.com/darklab8/fl-darkstat/configs/configs_mapped/freelancer_mapped/infocard_mapped/infocard"
	"github.com/darklab8/fl-darkstat/configs/configs_mapped/parserutils/filefind"
	"github.com/darklab8/fl-darkstat/configs/configs_mapped/parserutils/filefind/file"
	"github.com/darklab8/fl-darkstat/configs/configs_mapped/parserutils/iniload"
	"github.com/darklab8/fl-darkstat/configs/configs_mapped/parserutils/semantic"
	"github.com/darklab8/fl-darkstat/configs/configs_settings/logus"
	"github.com/darklab8/fl-darkstat/configs/overrides"
	"github.com/darklab8/fl-data-discovery/autopatcher"

	"github.com/darklab8/fl-darkstat/configs/configs_mapped/freelancer_mapped/data_mapped/equipment_mapped/equip_mapped"
	"github.com/darklab8/fl-darkstat/configs/configs_mapped/freelancer_mapped/data_mapped/equipment_mapped/market_mapped"
	"github.com/darklab8/fl-darkstat/configs/configs_mapped/freelancer_mapped/data_mapped/equipment_mapped/weaponmoddb"
	"github.com/darklab8/fl-darkstat/configs/discovery/base_recipe_items"
	"github.com/darklab8/fl-darkstat/configs/discovery/discoprices"
	"github.com/darklab8/fl-darkstat/configs/discovery/playercntl_rephacks"
	"github.com/darklab8/fl-darkstat/configs/discovery/pob_goods"
	"github.com/darklab8/fl-darkstat/configs/discovery/techcompat"

	"github.com/darklab8/go-utils/utils"
	"github.com/darklab8/go-utils/utils/timeit"
	"github.com/darklab8/go-utils/utils/utils_logus"
	"github.com/darklab8/go-utils/utils/utils_types"
)

type SiriusRevivalConfig struct {
	FLSRRecipes *flsr_recipes.Config
}

type DiscoveryConfig struct {
	Techcompat         *techcompat.Config
	Prices             *discoprices.Config
	BaseRecipeItems    *base_recipe_items.Config
	LatestPatch        autopatcher.Patch
	PlayercntlRephacks *playercntl_rephacks.Config
	PlayerOwnedBases   *pob_goods.Config
}

type MappedConfigs struct {
	FreelancerINI *exe_mapped.Config

	Universe *universe_mapped.Config
	Systems  *systems_mapped.Config

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

	NpcRankToDiff *npcranktodiff.Config
	DiffToMoney   *diff2money.Config

	FactionProps *faction_props_mapped.Config
	NpcShips     *npc_ships.Config
	Solararch    *solararch_mapped.Config
	Loadouts     *loadouts_mapped.Config

	Discovery *DiscoveryConfig
	FLSR      *SiriusRevivalConfig

	Overrides overrides.Overrides
}

func NewMappedConfigs() *MappedConfigs {
	return &MappedConfigs{}
}

func (configs *MappedConfigs) GetAvgTradeLaneSpeed() int {
	average_trade_lane_speed := 2250
	if configs.FLSR != nil {
		// make this value part of config files some day
		average_trade_lane_speed = 5000
	}
	return average_trade_lane_speed
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
	files_loadouts := getConfigs(filesystem, p.FreelancerINI.Loadouts)
	file_universe := iniload.NewLoader(filesystem.GetFile(universe_mapped.FILENAME))
	file_interface := iniload.NewLoader(filesystem.GetFile(interface_mapped.FILENAME_FL_INI))
	file_initialworld := iniload.NewLoader(filesystem.GetFile(initialworld.FILENAME))
	file_empathy := iniload.NewLoader(filesystem.GetFile(empathy_mapped.FILENAME))
	file_mbases := iniload.NewLoader(filesystem.GetFile(mbases_mapped.FILENAME))
	file_consts := iniload.NewLoader(filesystem.GetFile(const_mapped.FILENAME))
	file_weaponmoddb := iniload.NewLoader(filesystem.GetFile(weaponmoddb.FILENAME))

	file_diff2money := iniload.NewLoader(filesystem.GetFile(diff2money.FILENAME))
	file_npcranktodiff := iniload.NewLoader(filesystem.GetFile(npcranktodiff.FILENAME))

	file_faction_props := iniload.NewLoader(filesystem.GetFile(faction_props_mapped.FILENAME))
	file_npc_ships := iniload.NewLoader(filesystem.GetFile(npc_ships.FILENAME))
	file_solararch := iniload.NewLoader(filesystem.GetFile(solararch_mapped.FILENAME))

	all_files := append(files_goods, files_market...)
	all_files = append(all_files, files_equip...)
	all_files = append(all_files, files_shiparch...)
	all_files = append(all_files, files_loadouts...)
	all_files = append(all_files,
		file_universe,
		file_interface,
		file_initialworld,
		file_empathy,
		file_mbases,
		file_consts,
		file_weaponmoddb,
		file_diff2money,
		file_npcranktodiff,
		file_faction_props,
		file_npc_ships,
		file_solararch,
	)

	var file_techcompat *iniload.IniLoader
	var file_prices *iniload.IniLoader
	var file_base_recipe_items *iniload.IniLoader
	var file_playercntl_rephacks *iniload.IniLoader
	if filesystem.GetFile("flsr-launcher.ini") != nil ||
		filesystem.GetFile("flsr-texts.dll") != nil ||
		filesystem.GetFile("flsr-dialogs.dll") != nil {
		p.FLSR = &SiriusRevivalConfig{}
		flsr_recipes_file := filesystem.GetFile(flsr_recipes.FILENAME)
		if flsr_recipes_file != nil {
			p.FLSR.FLSRRecipes = flsr_recipes.Read(iniload.NewLoader(flsr_recipes_file).Scan())
		}
	}
	if techcom := filesystem.GetFile("launcherconfig.xml"); techcom != nil {
		p.Discovery = &DiscoveryConfig{}
		file_techcompat = iniload.NewLoader(file.NewWebFile("https://discoverygc.com/gameconfigpublic/techcompat.cfg"))
		file_prices = iniload.NewLoader(file.NewWebFile("https://discoverygc.com/gameconfigpublic/prices.cfg"))
		file_base_recipe_items = iniload.NewLoader(file.NewWebFile("https://discoverygc.com/gameconfigpublic/base_recipe_items.cfg"))
		file_playercntl_rephacks = iniload.NewLoader(file.NewWebFile("https://discoverygc.com/gameconfigpublic/playercntl_rephacks.cfg"))

		all_files = append(
			all_files,
			file_techcompat,
			file_prices,
			file_base_recipe_items,
			file_playercntl_rephacks,
		)

		if latest_patch_file := filesystem.GetFile(autopatcher.AutopatherFilename); latest_patch_file != nil {
			fmt.Println("latest_patch_file=", latest_patch_file)
			latest_patch_file_fp := latest_patch_file.GetFilepath()
			patch_data, err := os.ReadFile(latest_patch_file_fp.ToString())
			if !logus.Log.CheckError(err, "failed to unmarshal patch") {
				json.Unmarshal(patch_data, &p.Discovery.LatestPatch)
			}
			fmt.Println("p.Discovery.LatestPatch=", p.Discovery.LatestPatch)
		}
	}

	var infocards_override *file.File
	if p.Discovery != nil {
		infocards_override = file.NewWebFile("https://discoverygc.com/gameconfigpublic/infocard_overrides.cfg")
	}

	timeit.NewTimerF(func() {
		var wg sync.WaitGroup
		wg.Add(len(all_files))
		for _, file := range all_files {
			go func(file *iniload.IniLoader) {
				file.Scan()
				wg.Done()
			}(file)
		}
		wg.Wait()
	}, timeit.WithMsg("Scanned ini loaders"))

	overrides_file := filesystem.GetFile(overrides.FILENAME)
	if overrides_file != nil {
		logus.Log.Info("found overrides file")
		p.Overrides = overrides.Read(overrides_file.GetFilepath())
	}

	timeit.NewTimerF(func() {
		var wg sync.WaitGroup

		wg.Add(1)
		go func() {
			timeit.NewTimerF(func() {
				p.Universe = universe_mapped.Read(file_universe, filesystem)
				p.Systems = systems_mapped.Read(p.Universe, filesystem)
			}, timeit.WithMsg("map systems"))
			wg.Done()
		}()

		wg.Add(1)
		go func() {
			p.Market = market_mapped.Read(files_market)
			wg.Done()
		}()
		wg.Add(1)
		go func() {
			p.Equip = equip_mapped.Read(files_equip)
			wg.Done()
		}()
		wg.Add(1)
		go func() {
			p.Goods = equipment_mapped.Read(files_goods)
			wg.Done()
		}()
		wg.Add(1)
		go func() {
			p.Shiparch = ship_mapped.Read(files_shiparch)
			wg.Done()
		}()
		wg.Add(1)
		go func() {
			p.InfocardmapINI = interface_mapped.Read(file_interface)
			wg.Done()
		}()
		wg.Add(1)
		go func() {
			p.Infocards, _ = infocard_mapped.Read(filesystem, p.FreelancerINI, infocards_override)
			wg.Done()
		}()
		wg.Add(1)
		go func() {
			p.InitialWorld = initialworld.Read(file_initialworld)
			wg.Done()
		}()
		wg.Add(1)
		go func() {
			p.Empathy = empathy_mapped.Read(file_empathy)
			wg.Done()
		}()
		wg.Add(1)
		go func() {
			p.MBases = mbases_mapped.Read(file_mbases)
			wg.Done()
		}()
		wg.Add(1)
		go func() {
			p.Consts = const_mapped.Read(file_consts)
			wg.Done()
		}()
		wg.Add(1)
		go func() {
			p.WeaponMods = weaponmoddb.Read(file_weaponmoddb)
			wg.Done()
		}()
		wg.Add(1)
		go func() {
			p.NpcRankToDiff = npcranktodiff.Read(file_npcranktodiff)
			p.DiffToMoney = diff2money.Read(file_diff2money)
			wg.Done()
		}()

		wg.Add(1)
		go func() {
			p.FactionProps = faction_props_mapped.Read(file_faction_props)
			p.NpcShips = npc_ships.Read(file_npc_ships)
			wg.Done()
		}()

		wg.Add(1)
		go func() {
			p.Solararch = solararch_mapped.Read(file_solararch)
			wg.Done()
		}()

		wg.Add(1)
		go func() {
			p.Loadouts = loadouts_mapped.Read(files_loadouts)
			wg.Done()
		}()

		if p.Discovery != nil {
			wg.Add(4)
			go func() {
				p.Discovery.Techcompat = techcompat.Read(file_techcompat)
				wg.Done()
			}()
			go func() {
				p.Discovery.Prices = discoprices.Read(file_prices)
				wg.Done()
			}()
			go func() {
				p.Discovery.BaseRecipeItems = base_recipe_items.Read(file_base_recipe_items)
				wg.Done()
			}()
			go func() {
				p.Discovery.PlayercntlRephacks = playercntl_rephacks.Read(file_playercntl_rephacks)
				wg.Done()
			}()
			file_public_bases := file.NewWebFile("https://discoverygc.com/forums/base_admin.php?action=getjson")
			p.Discovery.PlayerOwnedBases = pob_goods.Read(file_public_bases)
		}
		wg.Wait()
	}, timeit.WithMsg("Mapped stuff"))

	logus.Log.Info("Parse OK for FreelancerFolderLocation=", utils_logus.FilePath(file1path))

	return p
}

type IsDruRun bool

func (p *MappedConfigs) Write(is_dry_run IsDruRun) {
	// write
	files := []*file.File{}

	files = append(files, p.Universe.Write())
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
