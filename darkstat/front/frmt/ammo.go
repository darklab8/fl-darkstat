package frmt

import (
	"fmt"

	"github.com/darklab8/fl-configs/configs/configs_export"
)

func GetAmmoLimitFormatted(ammo_limit configs_export.AmmoLimit) string {
	result := "" // fmt.Sprintf("%6.0f", gun.EnergyDamagePerSec)

	if ammo_limit.AmountInCatridge != nil {
		result += fmt.Sprintf("%6d", *ammo_limit.AmountInCatridge)
	}

	if ammo_limit.MaxCatridges != nil {
		result += fmt.Sprintf("(x%d)", *ammo_limit.MaxCatridges)
	}

	if result == "" {
		return "inf."
	}

	return result
}
