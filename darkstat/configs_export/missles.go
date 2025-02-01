package configs_export

type Missile struct {
	MaxAngularVelocity float64 `json:"max_angular_velocity"`
}

func (e *Exporter) GetMissiles(ids []*Tractor, buyable_ship_tech map[string]bool) []Gun {
	var missiles []Gun

	for _, gun_info := range e.Configs.Equip.Guns {
		missile, err := e.getGunInfo(gun_info, ids, buyable_ship_tech)

		if err != nil {
			continue
		}

		if missile.HpType == "" {
			continue
		}

		munition := e.Configs.Equip.MunitionMap[gun_info.ProjectileArchetype.Get()]
		if _, ok := munition.Motor.GetValue(); !ok {
			// Excluded regular guns
			continue
		}

		missile.MaxAngularVelocity, _ = munition.MaxAngularVelocity.GetValue()

		missiles = append(missiles, missile)
	}

	return missiles
}
