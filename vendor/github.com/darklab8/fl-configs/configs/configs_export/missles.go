package configs_export

type Missile struct {
	MaxAngularVelocity float64
}

func (e *Exporter) GetMissiles() []Gun {
	var missiles []Gun

	for _, gun_info := range e.configs.Equip.Guns {
		missile := e.getGunInfo(gun_info)

		if missile.HpType == "" {
			continue
		}

		munition := e.configs.Equip.MunitionMap[gun_info.ProjectileArchetype.Get()]
		if _, ok := munition.Motor.GetValue(); !ok {
			// Excluded regular guns
			continue
		}

		missile.MaxAngularVelocity, _ = munition.MaxAngularVelocity.GetValue()

		missiles = append(missiles, missile)
	}

	return missiles
}
