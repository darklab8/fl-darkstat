package configs_export

import (
	"math"
	"regexp"

	"github.com/darklab8/fl-darkstat/configs/cfg"
	"github.com/darklab8/fl-darkstat/configs/configs_mapped/freelancer_mapped/data_mapped/initialworld/flhash"
	"github.com/darklab8/fl-darkstat/configs/configs_settings/logus"
	"github.com/darklab8/fl-darkstat/darkstat/configs_export/infocarder"
	"github.com/darklab8/go-utils/typelog"
)

func (g Shield) GetTechCompat() *DiscoveryTechCompat { return g.DiscoveryTechCompat }

type Shield struct {
	Name string `json:"name" validate:"required"`

	Class      string `json:"class" validate:"required"`
	Type       string `json:"type" validate:"required"`
	Technology string `json:"technology" validate:"required"`
	Price      int    `json:"price" validate:"required"`

	Capacity          int     `json:"capacity" validate:"required"`
	RegenerationRate  int     `json:"regeneration_rate" validate:"required"`
	ConstantPowerDraw int     `json:"constant_power_draw" validate:"required"`
	Value             float64 `json:"value" validate:"required"`
	RebuildPowerDraw  int     `json:"rebuild_power_draw" validate:"required"`
	OffRebuildTime    int     `json:"off_rebuild_time" validate:"required"`

	Toughness float64 `json:"toughness" validate:"required"`
	HitPts    int     `json:"hit_pts" validate:"required"`
	Lootable  bool    `json:"lootable" validate:"required"`

	Nickname   string          `json:"nickname" validate:"required"`
	HpType     string          `json:"hp_type" validate:"required"`
	HpTypeHash flhash.HashCode `json:"-" swaggerignore:"true"`
	IdsName    int             `json:"ids_name" validate:"required"`
	IdsInfo    int             `json:"ids_info" validate:"required"`

	Bases map[cfg.BaseUniNick]*MarketGood `json:"-" swaggerignore:"true"`

	*DiscoveryTechCompat `json:"-" swaggerignore:"true"`
	Mass                 float64 `json:"mass" validate:"required"`
}

func (b Shield) GetNickname() string                       { return string(b.Nickname) }
func (b Shield) GetBases() map[cfg.BaseUniNick]*MarketGood { return b.Bases }

func (b Shield) GetDiscoveryTechCompat() *DiscoveryTechCompat { return b.DiscoveryTechCompat }

func (e *Exporter) GetShields(ids []*Tractor) []Shield {
	var shields []Shield

	for _, shield_gen := range e.Mapped.Equip().ShieldGens {
		shield := Shield{
			Nickname: shield_gen.Nickname.Get(),

			IdsInfo: shield_gen.IdsInfo.Get(),
			IdsName: shield_gen.IdsName.Get(),
			Bases:   make(map[cfg.BaseUniNick]*MarketGood),
		}
		shield.Mass, _ = shield_gen.Mass.GetValue()

		shield.Technology, _ = shield_gen.ShieldType.GetValue()
		shield.Capacity, _ = shield_gen.MaxCapacity.GetValue()

		shield.RegenerationRate, _ = shield_gen.RegenerationRate.GetValue()
		shield.ConstantPowerDraw, _ = shield_gen.ConstPowerDraw.GetValue()
		shield.RebuildPowerDraw, _ = shield_gen.RebuildPowerDraw.GetValue()
		shield.OffRebuildTime, _ = shield_gen.OfflineRebuildTime.GetValue()

		shield.Lootable, _ = shield_gen.Lootable.GetValue()
		shield.Toughness, _ = shield_gen.Toughness.GetValue()
		shield.HitPts, _ = shield_gen.HitPts.GetValue()

		if good_info, ok := e.Mapped.Goods.GoodsMap[shield.Nickname]; ok {
			if price, ok := good_info.Price.GetValue(); ok {
				shield.Price = price
				shield.Bases = e.GetAtBasesSold(GetCommodityAtBasesInput{
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

		e.exportInfocards(infocarder.InfocardKey(shield.Nickname), shield.IdsInfo)
		shield.DiscoveryTechCompat = CalculateTechCompat(e.Mapped.Discovery, ids, shield.Nickname)

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

func (e *Exporter) FilterToUsefulShields(shields []Shield) []Shield {
	var items []Shield = make([]Shield, 0, len(shields))
	for _, item := range shields {
		if !e.Buyable(item.Bases) {
			continue
		}
		items = append(items, item)
	}
	return items
}
