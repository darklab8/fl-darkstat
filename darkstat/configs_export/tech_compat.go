package configs_export

import (
	"sort"

	"github.com/darklab8/fl-darkstat/configs/cfg"
	"github.com/darklab8/fl-darkstat/configs/configs_mapped"
)

type TechCompatOrderer struct {
	cached_techcell_nil []CompatibleIDsForTractor
	configs             *configs_mapped.MappedConfigs
	exporter            *Exporter
}

func NewOrderedTechCompat(e *Exporter) *TechCompatOrderer {
	orderer := &TechCompatOrderer{
		configs:  e.Configs,
		exporter: e,
	}

	orderer.cached_techcell_nil = append(orderer.cached_techcell_nil, CompatibleIDsForTractor{
		TechCompat: e.Configs.Discovery.Techcompat.General.UnlistedTech.Get(),
		Tractor:    &Tractor{Name: "Most Factions"},
	})

	for _, faction := range e.Configs.Discovery.Techcompat.Factions {
		if unlisted_faction_modifier, ok := faction.DefaultUnlisted.GetValue(); ok {
			orderer.cached_techcell_nil = append(orderer.cached_techcell_nil, CompatibleIDsForTractor{
				TechCompat: unlisted_faction_modifier,
				Tractor:    orderer.exporter.TractorsByID[cfg.TractorID(faction.ID.Get())],
			})
		}
	}

	return orderer
}

func (orderer *TechCompatOrderer) GetOrederedTechCompat(DiscoveryTechCompat *DiscoveryTechCompat) []CompatibleIDsForTractor {
	var DiscoIDsCompatsOrdered []CompatibleIDsForTractor

	if DiscoveryTechCompat == nil {
		return DiscoIDsCompatsOrdered
	}

	if DiscoveryTechCompat.TechCell == "" {
		return orderer.cached_techcell_nil
	}

	for tractor_id, tech_tecompability := range DiscoveryTechCompat.TechcompatByID {
		if tech_tecompability < 11.0/100.0 {
			continue
		}

		if tractor, ok := orderer.exporter.TractorsByID[tractor_id]; ok {
			DiscoIDsCompatsOrdered = append(DiscoIDsCompatsOrdered, CompatibleIDsForTractor{
				TechCompat: tech_tecompability,
				Tractor:    tractor,
			})
		}
	}

	sort.Slice(DiscoIDsCompatsOrdered, func(i, j int) bool {
		return DiscoIDsCompatsOrdered[i].Tractor.Name < DiscoIDsCompatsOrdered[j].Tractor.Name
	})

	return DiscoIDsCompatsOrdered
}
