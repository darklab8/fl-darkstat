package lootable_shown

import (
	_ "embed"
	"strings"
)

//go:embed vanilla_wrecks.txt
var vanilla_wrecks string
var vanila_wrecks_cache map[string]bool = make(map[string]bool)

//go:embed vanilla_encounters.txt
var vanilla_encounters string
var vanila_encounters_cache map[string]bool = make(map[string]bool)

//go:embed flsr_wrecks.txt
var flsr_wrecks string
var flsr_wrecks_cache map[string]bool = make(map[string]bool)

//go:embed flsr_encounters.txt
var flsr_encounters string
var flsr_encounters_cache map[string]bool = make(map[string]bool)

//go:embed discovery_wrecks.txt
var discovery_wrecks string
var discovery_wrecks_cache map[string]bool = make(map[string]bool)

//go:embed discovery_encounters.txt
var discovery_encounters string
var discovery_encounters_cache map[string]bool = make(map[string]bool)

func GetGenericPermitted(data string, cached map[string]bool) map[string]bool {
	if len(cached) != 0 {
		return cached
	}

	cached = make(map[string]bool)
	for _, allowed_item := range strings.Split(data, "\n") {
		cached[allowed_item] = true
	}

	return cached
}

func GetVanillaWrecksAlloed() map[string]bool {
	return GetGenericPermitted(vanilla_wrecks, vanila_wrecks_cache)
}
func GetFLSRWrecksAllowed() map[string]bool {
	return GetGenericPermitted(flsr_wrecks, flsr_wrecks_cache)
}
func GetDiscoveryWrecksAllowed() map[string]bool {
	return GetGenericPermitted(discovery_wrecks, discovery_wrecks_cache)
}

func GetVanillaEncountersAlloed() map[string]bool {
	return GetGenericPermitted(vanilla_encounters, vanila_encounters_cache)
}
func GetFLSEncountersAllowed() map[string]bool {
	return GetGenericPermitted(flsr_encounters, flsr_encounters_cache)
}
func GetDiscoveryEncountersAllowed() map[string]bool {
	return GetGenericPermitted(discovery_encounters, discovery_encounters_cache)
}
