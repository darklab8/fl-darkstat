package configs_export

import (
	"math"
	"regexp"

	"github.com/darklab8/fl-configs/configs/configs_settings/logus"
	"github.com/darklab8/go-typelog/typelog"
)

type Shield struct {
	Name string

	Class      string
	Type       string
	Technology string
	Price      int

	Capacity          int
	RegenerationRate  int
	ConstantPowerDraw int
	Value             float64
	RebuildPowerDraw  int
	OffRebuildTime    int

	Toughness float64
	HitPts    int
	Lootable  bool

	Nickname string
	HpType   string
	IdsName  int
	IdsInfo  int

	Bases []GoodAtBase

	*DiscoveryTechCompat
}

func (e *Exporter) GetShields(ids []Tractor) []Shield {
	var shields []Shield

	for _, shield_gen := range e.configs.Equip.ShieldGens {
		shield := Shield{}

		shield.Nickname = shield_gen.Nickname.Get()

		if ids_info, ok := shield_gen.IdsInfo.GetValue(); ok {
			shield.IdsInfo = ids_info
		} else {
			continue
		}

		shield.IdsName = shield_gen.IdsName.Get()
		shield.Technology, _ = shield_gen.ShieldType.GetValue()

		if max_capacity, ok := shield_gen.MaxCapacity.GetValue(); ok {
			shield.Capacity = max_capacity
		} else {
			continue
		}

		shield.RegenerationRate = shield_gen.RegenerationRate.Get()
		shield.ConstantPowerDraw = shield_gen.ConstPowerDraw.Get()
		shield.RebuildPowerDraw = shield_gen.RebuildPowerDraw.Get()
		shield.OffRebuildTime = shield_gen.OfflineRebuildTime.Get()

		shield.Lootable, _ = shield_gen.Lootable.GetValue()
		shield.Toughness, _ = shield_gen.Toughness.GetValue()
		shield.HitPts = shield_gen.HitPts.Get()

		if good_info, ok := e.configs.Goods.GoodsMap[shield.Nickname]; ok {
			if price, ok := good_info.Price.GetValue(); ok {
				shield.Price = price
				shield.Bases = e.GetAtBasesSold(GetAtBasesInput{
					Nickname: good_info.Nickname.Get(),
					Price:    price,
				})

				var shield_value float64

				if shield.Capacity != 0 {
					shield_value = math.Abs(float64(shield.Capacity))
				} else if shield.RegenerationRate != 0 {
					shield_value = math.Abs(float64(shield.RegenerationRate))
				} else if shield.ConstantPowerDraw != 0 {
					shield_value = math.Abs(float64(shield.ConstantPowerDraw))
				}
				shield.Value = 1000 * shield_value / float64(shield.Price)
			}
		}

		shield.Name = e.GetInfocardName(shield.IdsName, shield.Nickname)

		if hp_type, ok := shield_gen.HpType.GetValue(); ok {
			shield.HpType = hp_type
			if parsed_type_class := TypeClassRegex.FindStringSubmatch(hp_type); len(parsed_type_class) > 0 {
				shield.Type = parsed_type_class[1]
				shield.Class = parsed_type_class[2]
			}
		}

		e.exportInfocards(InfocardKey(shield.Nickname), shield.IdsInfo)
		shield.DiscoveryTechCompat = CalculateTechCompat(e.configs.Discovery, ids, shield.Nickname)
		shields = append(shields, shield)
	}

	return shields
}

var TypeClassRegex *regexp.Regexp

func init() {
	TypeClassRegex = InitRegexExpression(`[a-zA-Z]+_([a-zA-Z]+)_[a-zA-Z_]+([0-9])`)
}

func InitRegexExpression(expression string) *regexp.Regexp {
	regex, err := regexp.Compile(string(expression))
	logus.Log.CheckPanic(err, "failed to init regex={%s} in ", typelog.String("expression", expression))
	return regex
}

func FilterToUsefulShields(shields []Shield) []Shield {
	var items []Shield = make([]Shield, 0, len(shields))
	for _, item := range shields {
		if !Buyable(item.Bases) {
			continue
		}
		items = append(items, item)
	}
	return items
}
