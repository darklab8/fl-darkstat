package configs_export

type Scanner struct {
	Name  string
	Price int

	Range          int
	CargoScanRange int

	Lootable bool
	Nickname string
	NameID   int
	InfoID   int

	Bases []GoodAtBase

	*DiscoveryTechCompat
}

func (e *Exporter) GetScanners(ids []Tractor) []Scanner {
	var scanners []Scanner

	for _, cm_info := range e.configs.Equip.Scanners {
		item := Scanner{}
		item.Nickname = cm_info.Nickname.Get()
		item.Lootable = cm_info.Lootable.Get()
		item.NameID = cm_info.IdsName.Get()
		item.InfoID = cm_info.IdsInfo.Get()

		if good_info, ok := e.configs.Goods.GoodsMap[item.Nickname]; ok {
			if price, ok := good_info.Price.GetValue(); ok {
				item.Price = price
				item.Bases = e.GetAtBasesSold(GetAtBasesInput{
					Nickname: good_info.Nickname.Get(),
					Price:    price,
				})
			}
		}

		if name, ok := e.configs.Infocards.Infonames[item.NameID]; ok {
			item.Name = string(name)
		}

		e.exportInfocards(InfocardKey(item.Nickname), item.InfoID)
		item.DiscoveryTechCompat = CalculateTechCompat(e.configs.Discovery, ids, item.Nickname)
		scanners = append(scanners, item)
	}
	return scanners
}

func FilterToUserfulScanners(items []Scanner) []Scanner {
	var useful_items []Scanner = make([]Scanner, 0, len(items))
	for _, item := range items {
		if !Buyable(item.Bases) {
			continue
		}
		useful_items = append(useful_items, item)
	}
	return useful_items
}
