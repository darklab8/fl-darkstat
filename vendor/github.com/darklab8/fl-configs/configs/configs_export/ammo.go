package configs_export

type Ammo struct {
	Name  string
	Price int

	HitPts    int
	AmmoLimit int
	Volume    int

	Nickname string
	NameID   int
	InfoID   int

	Bases []GoodAtBase

	*DiscoveryTechCompat
}

func (e *Exporter) GetAmmo(ids []Tractor) []Ammo {
	var tractors []Ammo

	for _, munition_info := range e.configs.Equip.Munitions {
		munition := Ammo{}
		munition.Nickname = munition_info.Nickname.Get()
		munition.NameID, _ = munition_info.IdsName.GetValue()
		munition.InfoID, _ = munition_info.IdsInfo.GetValue()

		if ammo_ids_name, ok := munition_info.IdsName.GetValue(); ok {
			if name, ok := e.configs.Infocards.Infonames[ammo_ids_name]; ok {
				munition.Name = string(name)
			} else {
				munition.Name = "undefined infoname"
			}
		} else {
			munition.Name = "undefined infoid"
		}

		munition.Price = -1
		if good_info, ok := e.configs.Goods.GoodsMap[munition_info.Nickname.Get()]; ok {
			if price, ok := good_info.Price.GetValue(); ok {
				munition.Price = price
				munition.Bases = e.GetAtBasesSold(GetAtBasesInput{
					Nickname:       good_info.Nickname.Get(),
					Price:          price,
					PricePerVolume: -1,
				})
			}
		}

		if !Buyable(munition.Bases) && (munition.Name == "undefined infoid" || munition.Name == "undefined infoname") {
			continue
		}

		e.exportInfocards(InfocardKey(munition.Nickname), munition.InfoID)
		munition.DiscoveryTechCompat = CalculateTechCompat(e.configs.Discovery, ids, munition.Nickname)
		tractors = append(tractors, munition)
	}
	return tractors
}

func FilterToUsefulAmmo(cms []Ammo) []Ammo {
	var useful_items []Ammo = make([]Ammo, 0, len(cms))
	for _, item := range cms {
		if !Buyable(item.Bases) {
			continue
		}
		useful_items = append(useful_items, item)
	}
	return useful_items
}
