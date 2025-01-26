package configs_mapped

func (configs *MappedConfigs) CraftableBaseName() string {
	if configs.Discovery != nil {
		return "PoB crafts"
	}
	if configs.FLSR != nil {
		return "Craftable"
	}

	return "NoCrafts"
}
