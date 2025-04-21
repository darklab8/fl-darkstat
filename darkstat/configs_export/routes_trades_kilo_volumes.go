package configs_export

import (
	"fmt"
	"math"
)

var (
	KiloVolume    float64 = 1000
	MaxKilVolumes float64 = 100
)

func KiloVolumesDeliverable(buying_good *MarketGood, selling_good *MarketGood) float64 {
	if buying_good.PoBGood == nil && selling_good.PoBGood == nil {
		return MaxKilVolumes
	}

	var buying_kilo_volume = MaxKilVolumes
	var selling_kilo_volume = MaxKilVolumes
	if buying_good.PoBGood != nil {
		if buying_good.PoBGood.Quantity <= buying_good.PoBGood.MinStock {
			buying_kilo_volume = 0
		} else {
			buying_kilo_volume = math.Min(MaxKilVolumes, (float64(buying_good.PoBGood.Quantity-buying_good.PoBGood.MinStock)*buying_good.Volume)/KiloVolume)
		}

		if buying_good.PoB.Money != nil {
			base_has_moner_for_vol := (buying_good.Volume * float64(*buying_good.PoB.Money) / float64(buying_good.PriceBaseSellsFor) / KiloVolume)
			buying_kilo_volume = math.Min(base_has_moner_for_vol, buying_kilo_volume)
		}
	}

	if selling_good.PoBGood != nil {
		if selling_good.PoB.Nickname == "cayman_research_institute" && selling_good.Nickname == "commodity_high_temp_alloys" {
			fmt.Println()
		}

		if selling_good.PoBGood.Quantity >= selling_good.PoBGood.MaxStock {
			selling_kilo_volume = 0
		} else {
			selling_kilo_volume = math.Min(MaxKilVolumes, (float64(selling_good.PoBGood.MaxStock-selling_good.PoBGood.Quantity)*selling_good.Volume)/KiloVolume)
		}

		if selling_good.PoB.Money != nil && selling_good.PriceBaseBuysFor != nil {
			base_has_moner_for_vol := (selling_good.Volume * float64(*selling_good.PoB.Money) / float64(*selling_good.PriceBaseBuysFor)) / KiloVolume
			selling_kilo_volume = math.Min(base_has_moner_for_vol, selling_kilo_volume)
		}

	}

	return math.Min(MaxKilVolumes, math.Min(buying_kilo_volume, selling_kilo_volume))
}
