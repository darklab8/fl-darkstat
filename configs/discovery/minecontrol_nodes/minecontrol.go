package minecontrol_nodes

import (
	"github.com/darklab8/fl-darkstat/configs/cfg"
	"github.com/darklab8/fl-darkstat/configs/configs_mapped/parserutils/iniload"
	"github.com/darklab8/fl-darkstat/configs/configs_mapped/parserutils/semantic"
)

type NodeArchetype struct {
	semantic.Model
	Nickname *semantic.String
}

type MiningSystem struct {
	semantic.Model
	Nickname        *semantic.String
	NodeArchetypes  []*NodeArchetype
	Position        cfg.Vector
	SystemNickname  *semantic.String
	MaxSpawnCount   *semantic.Int
	RespawnCooldown *semantic.Int
}

type Config struct {
	MiningSystems             []*MiningSystem
	MiningSystemsBySystemNick map[string][]*MiningSystem
}

func Read(input_file *iniload.IniLoader) *Config {
	conf := &Config{
		MiningSystemsBySystemNick: map[string][]*MiningSystem{},
	}

	if resources, ok := input_file.SectionMap["[miningsystem]"]; ok {
		for _, resource := range resources {
			mining_system := &MiningSystem{
				Nickname:        semantic.NewString(resource, "name"),
				SystemNickname:  semantic.NewString(resource, "system"),
				MaxSpawnCount:   semantic.NewInt(resource, "max_spawn_count"),
				RespawnCooldown: semantic.NewInt(resource, "respawn_cooldown"),
			}
			mining_system.Map(resource)

			for index, param := range resource.ParamMap["position"] {
				position := semantic.NewVector(resource, param.Key, semantic.Precision(0), semantic.Index(index))
				mining_system.Position.X += position.Get().X
				mining_system.Position.Y += position.Get().Y
				mining_system.Position.Z += position.Get().Z
			}

			points_amount := len(resource.ParamMap["position"])
			mining_system.Position.X /= float64(points_amount)
			mining_system.Position.Y /= float64(points_amount)
			mining_system.Position.Z /= float64(points_amount)

			for index, param := range resource.ParamMap["node_archetype"] {

				node_archetype := &NodeArchetype{
					Nickname: semantic.NewString(resource, param.Key, semantic.OptsS(semantic.Index(index))),
				}
				node_archetype.Map(resource)
				mining_system.NodeArchetypes = append(mining_system.NodeArchetypes, node_archetype)
			}

			conf.MiningSystems = append(conf.MiningSystems, mining_system)
			conf.MiningSystemsBySystemNick[mining_system.SystemNickname.Get()] = append(conf.MiningSystemsBySystemNick[mining_system.SystemNickname.Get()], mining_system)

		}
	}

	return conf
}
