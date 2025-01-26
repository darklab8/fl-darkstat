package configs_export

import (
	"errors"
	"math"

	"github.com/darklab8/fl-darkstat/configs/configs_mapped/freelancer_mapped/data_mapped/universe_mapped/systems_mapped"
	"github.com/darklab8/fl-darkstat/configs/configs_mapped/parserutils/semantic"
	"github.com/darklab8/fl-darkstat/configs/configs_settings/logus"
)

func DistanceForVecs(Pos1 *semantic.Vect, Pos2 *semantic.Vect) (float64, error) {
	if _, ok := Pos1.X.GetValue(); !ok {
		return 0, errors.New("no x")
	}
	if _, ok := Pos2.X.GetValue(); !ok {
		return 0, errors.New("no x")
	}

	x_dist := math.Pow((Pos1.X.Get() - Pos2.X.Get()), 2)
	y_dist := math.Pow((Pos1.Y.Get() - Pos2.Y.Get()), 2)
	z_dist := math.Pow((Pos1.Z.Get() - Pos2.Z.Get()), 2)
	distance := math.Pow((x_dist + y_dist + z_dist), 0.5)
	return distance, nil
}

func GetMaxRadius(Size *semantic.Vect) (float64, error) {
	max_size := 0.0
	if value, ok := Size.X.GetValue(); ok {
		if value > max_size {
			max_size = value
		}
	}
	if value, ok := Size.Y.GetValue(); ok {
		if value > max_size {
			max_size = value
		}
	}
	if value, ok := Size.Z.GetValue(); ok {
		if value > max_size {
			max_size = value
		}
	}
	if max_size == 0 {
		return 0, errors.New("not found size")
	}

	return max_size, nil
}

func IsAnyVignetteWithinNPCSpawnRange(system *systems_mapped.System, npc_spawn_zone *systems_mapped.MissionPatrolZone) bool {
	matched_vignette := false
	for _, vignette := range system.MissionZoneVignettes {

		distance, dist_err := DistanceForVecs(vignette.Pos, npc_spawn_zone.Pos)
		if dist_err != nil {
			continue
		}

		max_spwn_zone_size, err_max_size := GetMaxRadius(npc_spawn_zone.Size)
		logus.Log.CheckWarn(err_max_size, "expected finding max size, but object does not have it")

		if distance < float64(vignette.Size.Get())+max_spwn_zone_size {
			matched_vignette = true
			break
		}
	}

	return matched_vignette
}
