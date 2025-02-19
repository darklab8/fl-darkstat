package darkgrpc

import (
	"context"

	pb "github.com/darklab8/fl-darkstat/darkapis/darkgrpc/statproto"
	"github.com/darklab8/fl-darkstat/darkstat/configs_export"
	"github.com/darklab8/go-utils/utils/ptr"
)

func (s *Server) GetPoBs(_ context.Context, in *pb.Empty) (*pb.GetPoBsReply, error) {
	if s.app_data != nil {
		s.app_data.Lock()
		defer s.app_data.Unlock()
	}

	var bases []*pb.PoB
	for _, base := range s.app_data.Configs.PoBs {

		item := &pb.PoB{
			Core: NewPoBCore(&base.PoBCore),
		}
		for _, shop_item := range base.ShopItems {
			item.ShopItems = append(item.ShopItems, NewShopItem(shop_item))
		}

		bases = append(bases, item)
	}
	return &pb.GetPoBsReply{Items: bases}, nil
}

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

func NewInt64FromDefenseMode(DefenseMode *configs_export.DefenseMode) *int64 {
	if DefenseMode == nil {
		return nil
	}

	return ptr.Ptr(int64(*DefenseMode))
}

func NewShopItem(item *configs_export.ShopItem) *pb.ShopItem {
	return &pb.ShopItem{
		Nickname:  item.Nickname,
		Name:      item.Name,
		Category:  item.Category,
		Id:        int64(item.Id),
		Quantity:  int64(item.Quantity),
		Price:     int64(item.Price),
		SellPrice: int64(item.SellPrice),
		MinStock:  int64(item.MinStock),
		MaxStock:  int64(item.MaxStock),
	}
}

func (s *Server) GetPoBGoods(_ context.Context, in *pb.Empty) (*pb.GetPoBGoodsReply, error) {
	if s.app_data != nil {
		s.app_data.Lock()
		defer s.app_data.Unlock()
	}

	var pob_goods []*pb.PoBGood
	for _, base := range s.app_data.Configs.PoBGoods {

		item := &pb.PoBGood{
			Nickname:              base.Nickname,
			Name:                  base.Name,
			TotalBuyableFromBases: int64(base.TotalBuyableFromBases),
			TotalSellableToBases:  int64(base.TotalSellableToBases),
			BestPriceToBuy:        NewInt64(base.BestPriceToBuy),
			BestPriceToSell:       NewInt64(base.BestPriceToSell),
			Category:              base.Category,
			AnyBaseSells:          base.AnyBaseSells,
			AnyBaseBuys:           base.AnyBaseBuys,
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
