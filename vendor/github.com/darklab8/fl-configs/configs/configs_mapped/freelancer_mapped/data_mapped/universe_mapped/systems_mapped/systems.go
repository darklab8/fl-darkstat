package systems_mapped

import (
	"fmt"

	"github.com/darklab8/fl-configs/configs/configs_mapped/freelancer_mapped/data_mapped/universe_mapped"
	"github.com/darklab8/fl-configs/configs/configs_mapped/parserutils/filefind"
	"github.com/darklab8/fl-configs/configs/configs_mapped/parserutils/filefind/file"
	"github.com/darklab8/fl-configs/configs/configs_mapped/parserutils/inireader"
	"github.com/darklab8/fl-configs/configs/configs_mapped/parserutils/semantic"
	"github.com/darklab8/go-utils/goutils/utils/utils_types"
)

const (
	KEY_OBJECT   = "[Object]"
	KEY_NICKNAME = "nickname"
	KEY_BASE     = "base"
)

type Base struct {
	semantic.Model
	Nickname *semantic.String
	Base     *semantic.String // base.nickname in universe.ini
	DockWith *semantic.String

	IDsInfo     *semantic.Int
	IdsName     *semantic.Int
	RepNickname *semantic.String
	Pos         *semantic.Vect
}
type System struct {
	semantic.ConfigModel
	Nickname     string
	Bases        []*Base
	BasesByNick  map[string]*Base
	BasesByBases map[string][]*Base
}

type Config struct {
	SystemsMap map[string]*System
	Systems    []*System
}

type FileRead struct {
	system_key string
	file       *file.File
	ini        *inireader.INIFile
}

func Read(universe_config *universe_mapped.Config, filesystem *filefind.Filesystem) *Config {
	frelconfig := &Config{}

	var system_files map[string]*file.File = make(map[string]*file.File)
	for _, base := range universe_config.Bases {
		base_system := base.System.Get()
		fmt.Println("base_system=", base_system)
		if base_system == "ew06" {
			fmt.Println()
		}
		universe_system := universe_config.SystemMap[universe_mapped.SystemNickname(base_system)]

		if universe_system == nil {
			fmt.Println("base_system==", base_system)
			for key, _ := range universe_config.SystemMap {
				fmt.Print("key=", key)
			}
			fmt.Println()
		}

		var filename utils_types.FilePath
		// func() {
		// 	defer func() {
		// 		if r := recover(); r != nil {
		// 			logus.Log.Error("Recovered from int File.FileName Error:\n",
		// 				typelog.Any("recover", r),
		// 				typelog.Any("universe_system", universe_system),
		// 				typelog.Any("base_system", base_system),
		// 				// typelog.Any("universe.File", *(universe_system.File)),
		// 			)
		// 			panic(r)
		// 		}
		// 	}()
		//
		// }()
		filename = universe_system.File.FileName()
		path := filesystem.GetFile(filename)
		if path.GetFilepath() == "" {
			fmt.Println()
		}
		system_files[base.System.Get()] = file.NewFile(path.GetFilepath())
	}

	var system_iniconfigs map[string]*inireader.INIFile = make(map[string]*inireader.INIFile)

	func() {
		// Read system files with parallelism ^_^
		iniconfigs_channel := make(chan *FileRead)
		read_file := func(data *FileRead) {
			data.ini = inireader.Read(data.file)
			iniconfigs_channel <- data
		}
		for system_key, file := range system_files {
			go read_file(&FileRead{
				system_key: system_key,
				file:       file,
			})
		}
		for range system_files {
			result := <-iniconfigs_channel
			system_iniconfigs[result.system_key] = result.ini
		}
	}()

	frelconfig.SystemsMap = make(map[string]*System)
	frelconfig.Systems = make([]*System, 0)
	for system_key, sysiniconf := range system_iniconfigs {
		system_to_add := &System{}
		system_to_add.Init(sysiniconf.Sections, sysiniconf.Comments, sysiniconf.File.GetFilepath())

		system_to_add.Nickname = system_key
		system_to_add.BasesByNick = make(map[string]*Base)
		system_to_add.BasesByBases = make(map[string][]*Base)
		system_to_add.Bases = make([]*Base, 0)
		frelconfig.SystemsMap[system_key] = system_to_add
		frelconfig.Systems = append(frelconfig.Systems, system_to_add)

		if objects, ok := sysiniconf.SectionMap[KEY_OBJECT]; ok {
			for _, obj := range objects {

				// check if it is base object
				_, ok := obj.ParamMap[KEY_BASE]
				if ok {
					base_to_add := &Base{}
					base_to_add.Map(obj)

					base_to_add.Nickname = semantic.NewString(obj, KEY_NICKNAME, semantic.WithLowercaseS(), semantic.WithoutSpacesS())
					base_to_add.Base = semantic.NewString(obj, KEY_BASE, semantic.WithLowercaseS(), semantic.WithoutSpacesS())
					base_to_add.DockWith = semantic.NewString(obj, "dock_with", semantic.OptsS(semantic.Optional()))
					base_to_add.RepNickname = semantic.NewString(obj, "reputation", semantic.OptsS(semantic.Optional()), semantic.WithLowercaseS(), semantic.WithoutSpacesS())

					base_to_add.IDsInfo = semantic.NewInt(obj, "ids_info", semantic.Optional())
					base_to_add.IdsName = semantic.NewInt(obj, "ids_name", semantic.Optional())

					base_to_add.Pos = semantic.NewVector(obj, "pos", semantic.Precision(0))

					system_to_add.BasesByNick[base_to_add.Nickname.Get()] = base_to_add
					system_to_add.BasesByBases[base_to_add.Base.Get()] = append(system_to_add.BasesByBases[base_to_add.Base.Get()], base_to_add)
					system_to_add.Bases = append(system_to_add.Bases, base_to_add)
				}
			}
		}

	}

	return frelconfig
}

func (frelconfig *Config) Write() []*file.File {
	var files []*file.File = make([]*file.File, 0)
	for _, system := range frelconfig.Systems {
		inifile := system.Render()
		files = append(files, inifile.Write(inifile.File))
	}
	return files
}
