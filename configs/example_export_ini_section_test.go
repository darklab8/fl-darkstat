package configs

import (
	"fmt"
	"strings"

	"github.com/darklab8/fl-darkstat/configs/configs_mapped"
	"github.com/darklab8/fl-darkstat/configs/configs_mapped/parserutils/inireader"
	"github.com/darklab8/fl-darkstat/configs/configs_settings"
	"github.com/darklab8/fl-darkstat/configs/configs_settings/logus"
	"github.com/darklab8/go-utils/utils/utils_logus"
)

// ExampleExtractIniSection demonstrates how to extract specific freelancer ini sections out
func Example_extractIniSection() {

	freelancer_folder := configs_settings.Env.FreelancerFolder
	configs := configs_mapped.NewMappedConfigs()
	logus.Log.Debug("scanning freelancer folder", utils_logus.FilePath(freelancer_folder))

	// Reading to ini universal custom format and mapping to ORM objects
	// which have both reading and writing back capabilities
	configs.Read(freelancer_folder)

	order_gun := configs.Equip.GunMap["fc_or_gun01_mark02"]
	var output strings.Builder

	order_gun_section := order_gun.Model.RenderModel()
	output.WriteString(fmt.Sprintf("%s\n", string(order_gun_section.OriginalType)))

	for _, param := range order_gun_section.Params {
		output.WriteString(fmt.Sprintf("%s\n", param.ToString(inireader.WithComments(true))))
	}

	order_munition := configs.Equip.MunitionMap[order_gun.ProjectileArchetype.Get()]

	order_munition_section := order_munition.Model.RenderModel()
	output.WriteString(fmt.Sprintf("%s\n", string(order_munition_section.OriginalType)))

	for _, param := range order_munition_section.Params {
		output.WriteString(fmt.Sprintf("%s\n", param.ToString(inireader.WithComments(true))))
	}

	fmt.Println(output.String())

	// Example of output
	// [Gun]
	// nickname = fc_or_gun01_mark02
	// ids_name = 263221
	// ids_info = 264221
	// da_archetype = equipment\models\weapons\li_heavy_ion_blaster.cmp
	// material_library = equipment\models\li_equip.mat
	// hp_child = HPConnect
	// hit_pts = 7500
	// explosion_resistance = 0.14
	// debris_type = debris_normal
	// parent_impulse = 20
	// child_impulse = 0.08
	// volume = 0
	// mass = 0.01
	// hp_gun_type = hp_gun_special_1
	// damage_per_fire = 0
	// power_usage = 56.9
	// refire_delay = 0.120000
	// muzzle_velocity = 800
	// toughness = 4.200000
	// flash_particle_name = pi_death_01_flash
	// flash_radius = 15
	// light_anim = l_gun01_flash
	// projectile_archetype = fc_or_gun01_mark02_ammo
	// separation_explosion = sever_debris
	// auto_turret = false
	// turn_rate = 90
	// lootable = true
	// lodranges = 0, 40, 70, 100, 150
	// [Munition]
	// nickname = fc_or_gun01_mark02_ammo
	// hp_type = hp_gun
	// requires_ammo = false
	// hit_pts = 2
	// hull_damage = 84
	// energy_damage = 21
	// armor_pen = 0.2
	// weapon_type = W_Laser01
	// one_shot_sound = fire_laser5
	// munition_hit_effect = pi_laser_04_impact
	// const_effect = pi_death_01_proj
	// lifetime = 0.7
	// force_gun_ori = false
	// mass = 1
	// volume = 0
}
