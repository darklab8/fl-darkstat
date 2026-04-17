package configs_export

import (
	"math"

	"github.com/darklab8/fl-darkstat/configs/configs_mapped/freelancer_mapped/data_mapped/equipment_mapped/equip_mapped"
)

type Missile struct {
	MaxAngularVelocity float64 `json:"max_angular_velocity"  validate:"required"`
	TopSpeedAtRest     float64 `json:"top_speed_at_rest"  validate:"required"`
	TopSpeedAt80       float64 `json:"top_speed_at_80"  validate:"required"`
	TopSpeedAt200      float64 `json:"top_speed_at_200"  validate:"required"`

	RangeAt80Speed  float64 `json:"range_at_80_speed"  validate:"required"`
	RangeAt200Speed float64 `json:"range_at_200_speed"  validate:"required"`
}

// GetRangeAtSpeed calculates the effective range of a munition given motor and velocity parameters.
// Returns distance in meters as a float64.
func GetRangeAtSpeed(
	motor *equip_mapped.Motor,
	munition *equip_mapped.Munition,
	muzzleVelocity float64,
	shipVelocity float64,
	topSpeed float64,
) float64 {
	startSpeed := muzzleVelocity + shipVelocity
	distance := munition.LifeTime.Get() * startSpeed

	if motor != nil && motor.Acceleration.Get() > 0.0 {
		motorLifetime := min(motor.Lifetime.Get(), munition.LifeTime.Get()-motor.Delay.Get())

		startSpeed = min(startSpeed, topSpeed)
		motorLifetime = min(motorLifetime, max(0.0, topSpeed-startSpeed)/motor.Acceleration.Get())

		motorAcceleration := motor.Acceleration.Get() * motorLifetime

		distance += motorAcceleration * motorLifetime * 0.5

		timeMaxSpeed := max(0.0, munition.LifeTime.Get()-motor.Delay.Get()-motorLifetime)
		distance += motorAcceleration * timeMaxSpeed
	}

	return distance
}

func getMunitionTopSpeed(
	motor *equip_mapped.Motor,
	munition *equip_mapped.Munition,
	muzzleVelocity float64,
) (TopSpeed float64, Motorlifetime float64) {
	if motor == nil || motor.Acceleration.Get() <= 0 {
		TopSpeed = muzzleVelocity
		return
	}
	if value, _ := motor.Acceleration.GetValue(); value <= 0 {
		TopSpeed = muzzleVelocity
		return
	}

	Motorlifetime = min(motor.Lifetime.Get(), munition.LifeTime.Get()-motor.Delay.Get())
	TopSpeed = muzzleVelocity + motor.Acceleration.Get()*Motorlifetime
	return
}

func (e *Exporter) GetMissiles(ids []*Tractor, buyable_ship_tech map[string]bool) []Gun {
	var missiles []Gun

	for _, gun_info := range e.Mapped.Equip().Guns {
		missile, err := e.getGunInfo(gun_info, ids, buyable_ship_tech)

		if err != nil {
			continue
		}

		if missile.HpType == "" {
			continue
		}

		munition := e.Mapped.Equip().MunitionMap[gun_info.ProjectileArchetype.Get()]
		if _, ok := munition.Motor.GetValue(); !ok {
			// Excluded regular guns
			continue
		}

		var motor *equip_mapped.Motor
		if motor_nick, ok := munition.Motor.GetValue(); ok {
			motor = e.Mapped.Equip().MotorMap[motor_nick]
		}

		missile.MaxAngularVelocity, _ = munition.MaxAngularVelocity.GetValue()

		missile.TopSpeedAtRest, _ = getMunitionTopSpeed(motor, munition, gun_info.MuzzleVelosity.Get())
		missile.TopSpeedAt80 = missile.TopSpeedAtRest + 80
		missile.TopSpeedAt200 = missile.TopSpeedAtRest + 200

		if e.Mapped.Discovery != nil {
			if value, ok := munition.TopSpeed.GetValue(); ok {
				missile.TopSpeedAtRest = math.Min(missile.TopSpeedAtRest, value)
				missile.TopSpeedAt80 = math.Min(missile.TopSpeedAt80, value)
				missile.TopSpeedAt200 = math.Min(missile.TopSpeedAt200, value)
			}
		}

		missile.Range = GetRangeAtSpeed(motor, munition, gun_info.MuzzleVelosity.Get(), 0, missile.TopSpeedAtRest)
		missile.RangeAt80Speed = GetRangeAtSpeed(motor, munition, gun_info.MuzzleVelosity.Get(), 80, missile.TopSpeedAt80)
		missile.RangeAt200Speed = GetRangeAtSpeed(motor, munition, gun_info.MuzzleVelosity.Get(), 200, missile.TopSpeedAt200)

		missiles = append(missiles, missile)
	}

	return missiles
}
