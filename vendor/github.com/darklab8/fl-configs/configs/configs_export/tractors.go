package configs_export

type Tractor struct {
	Name       string
	Price      int
	MaxLength  int
	ReachSpeed int

	Lootable bool
	Nickname string
	NameID   int
	InfoID   int

	Bases []GoodAtBase
}

func (e *Exporter) GetTractors() []Tractor {
	var tractors []Tractor

	for _, tractor_info := range e.configs.Equip.Tractors {
		tractor := Tractor{}
		tractor.Nickname = tractor_info.Nickname.Get()
		tractor.MaxLength = tractor_info.MaxLength.Get()
		tractor.ReachSpeed = tractor_info.ReachSpeed.Get()
		tractor.Lootable = tractor_info.Lootable.Get()
		tractor.NameID = tractor_info.IdsName.Get()
		tractor.InfoID = tractor_info.IdsInfo.Get()

		if good_info, ok := e.configs.Goods.GoodsMap[tractor.Nickname]; ok {
			if price, ok := good_info.Price.GetValue(); ok {
				tractor.Price = price
				tractor.Bases = e.GetAtBasesSold(GetAtBasesInput{
					Nickname:       good_info.Nickname.Get(),
					Price:          price,
					PricePerVolume: -1,
				})
			}
		}

		if name, ok := e.configs.Infocards.Infonames[tractor.NameID]; ok {
			tractor.Name = string(name)
		}

		e.exportInfocards(InfocardKey(tractor.Nickname), tractor.InfoID)
		tractors = append(tractors, tractor)
	}
	return tractors
}

func FilterToUsefulTractors(tractors []Tractor) []Tractor {
	var buyable_tractors []Tractor = make([]Tractor, 0, len(tractors))
	for _, item := range tractors {
		if !Buyable(item.Bases) {
			continue
		}
		buyable_tractors = append(buyable_tractors, item)
	}
	return buyable_tractors
}
