package darkgrpc_deprecated

import (
	"net/http"
	"strconv"

	pb "github.com/darklab8/fl-darkstat/darkapis/darkgrpc_deprecated/statproto_deprecated"
	"github.com/darklab8/fl-darkstat/darkapis/darkhttp/apiutils"
	"github.com/darklab8/fl-darkstat/darkcore/web"
	"github.com/darklab8/fl-darkstat/darkcore/web/registry"
	"github.com/darklab8/fl-darkstat/darkstat/appdata"
	"github.com/darklab8/fl-darkstat/darkstat/configs_export"
	"github.com/darklab8/go-utils/utils/ptr"
)

func NewPoBCore(base *configs_export.PoBCore) *pb.PoBCore {
	return &pb.PoBCore{
		Nickname: base.Nickname,
		Name:     base.Name,

		Pos:         base.Pos,
		Level:       NewInt64(base.Level),
		Money:       NewInt64(base.Money),
		Health:      base.Health,
		DefenseMode: NewInt64FromDefenseMode(base.DefenseMode),

		SystemNick:  base.SystemNick,
		SystemName:  base.SystemName,
		FactionNick: base.FactionNick,
		FactionName: base.FactionName,

		ForumThreadUrl: base.ForumThreadUrl,
		CargoSpaceLeft: NewInt64(base.CargoSpaceLeft),

		BasePos:     NewPos(base.BasePos),
		SectorCoord: base.SectorCoord,
		Region:      base.Region,
	}
}

func NewInt64FromDefenseMode(DefenseMode *configs_export.DefenseMode) pb.NumString {
	if DefenseMode == nil {
		return nil
	}

	return ptr.Ptr(strconv.Itoa(int(*DefenseMode)))
}

func NewShopItem(item *configs_export.ShopItem) *pb.ShopItem {
	return &pb.ShopItem{
		Nickname:  item.Nickname,
		Name:      item.Name,
		Category:  item.Category,
		Id:        NewInt64S(item.Id),
		Quantity:  NewInt64S(item.Quantity),
		Price:     NewInt64S(item.PriceBaseSellsFor),
		SellPrice: NewInt64S(item.PriceBaseBuysFor),
		MinStock:  NewInt64S(item.MinStock),
		MaxStock:  NewInt64S(item.MaxStock),
	}
}

func GetPoBGoods(app_data *appdata.AppData) (*pb.GetPoBGoodsReply, error) {
	var pob_goods []*pb.PoBGood
	for _, base := range app_data.Configs.PoBGoods {

		item := &pb.PoBGood{
			Nickname:              base.Nickname,
			Name:                  base.Name,
			TotalBuyableFromBases: NewInt64S(base.TotalBuyableFromBases),
			TotalSellableToBases:  NewInt64S(base.TotalSellableToBases),
			BestPriceToBuy:        NewInt64(base.BestPriceToBuy),
			BestPriceToSell:       NewInt64(base.BestPriceToSell),
			Category:              base.Category,
			AnyBaseSells:          base.AnyBaseSells,
			AnyBaseBuys:           base.AnyBaseBuys,
			Volume:                base.Volume,
			ShipClass:             NewShipClass(base.ShipClass),
		}

		for _, shop_item := range base.Bases {
			item.Bases = append(item.Bases, &pb.PoBGoodBase{
				ShopItem: NewShopItem(shop_item.ShopItem),
				Base:     NewPoBCore(shop_item.Base),
			})
		}

		pob_goods = append(pob_goods, item)
	}
	return &pb.GetPoBGoodsReply{Items: pob_goods}, nil
}

type Api interface {
	GetAppData() *appdata.AppData
}

// ShowAccount godoc
// @Summary      PoB Goods Deprecated
// @Tags         deprecated
// @Accept       json
// @Produce      json
// @Success      200  {object}  	statproto_deprecated.GetPoBGoodsReply
// @Router       /statproto.Darkstat/GetPoBGoods [post]
func GetPobGoodsDeprecated(webapp *web.Web, api Api) *registry.Endpoint {
	return &registry.Endpoint{
		Url: "" + "/statproto.Darkstat/GetPoBGoods",
		Handler: func(w http.ResponseWriter, r *http.Request) {
			if webapp.AppDataMutex != nil {
				webapp.AppDataMutex.RLock()
				defer webapp.AppDataMutex.RUnlock()
			}

			reply, err := GetPoBGoods(api.GetAppData())

			if err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}

			apiutils.ReturnJson(&w, reply)
		},
	}
}
