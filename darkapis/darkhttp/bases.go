package darkhttp

import (
	"net/http"

	pb "github.com/darklab8/fl-darkstat/darkapis/darkgrpc/statproto"
	"github.com/darklab8/fl-darkstat/darkapis/darkhttp/apiutils"
	"github.com/darklab8/fl-darkstat/darkapis/services"
	"github.com/darklab8/fl-darkstat/darkcore/web"
	"github.com/darklab8/fl-darkstat/darkcore/web/registry"
	"github.com/darklab8/fl-darkstat/darkstat/configs_export"
)

type Base struct {
	*configs_export.Base
	MarketGoods []*configs_export.MarketGood `json:"market_goods"`
}

// ShowAccount godoc
// @Summary      Getting list of NPC Bases
// @Tags         bases
// @Accept       json
// @Produce      json
// @Success      200  {array}  	darkapi.Base
// @Router       /api/npc_bases [post]
// @Param request body pb.GetBasesInput true "input variables"
func GetBases(webapp *web.Web, api *Api) *registry.Endpoint {
	return &registry.Endpoint{
		Url: "GET " + ApiRoute + "/npc_bases",
		Handler: func(w http.ResponseWriter, r *http.Request) {
			if webapp.AppDataMutex != nil {
				webapp.AppDataMutex.Lock()
				defer webapp.AppDataMutex.Unlock()
			}

			var in pb.GetBasesInput
			ReadJsonInput(w, r, &in)
			filter_to_useful := in.FilterToUseful
			include_market_goods := in.IncludeMarketGoods

			var result []*configs_export.Base
			if filter_to_useful {
				result = configs_export.FilterToUserfulBases(api.app_data.Configs.Bases)
			} else {
				result = api.app_data.Configs.Bases
			}
			result = services.FilterNicknames(in.FilterNicknames, result)

			var output []*Base
			for _, item := range result {
				answer := &Base{
					Base: item,
				}
				if include_market_goods {
					answer.MarketGoods = services.FilterMarketGoodCategory(in.FilterMarketGoodCategory, item.MarketGoodsPerNick)
				}
				output = append(output, answer)
			}
			apiutils.ReturnJson(&w, output)
		},
	}

}

// ShowAccount godoc
// @Summary      Getting list of Mining Operations
// @Tags         bases
// @Accept       json
// @Produce      json
// @Success      200  {array}  	darkapi.Base
// @Router       /api/mining_operations [post]
// @Param request body pb.GetBasesInput true "input variables"
func GetOreFields(webapp *web.Web, api *Api) *registry.Endpoint {
	return &registry.Endpoint{
		Url: "GET " + ApiRoute + "/mining_operations",
		Handler: func(w http.ResponseWriter, r *http.Request) {
			if webapp.AppDataMutex != nil {
				webapp.AppDataMutex.Lock()
				defer webapp.AppDataMutex.Unlock()
			}
			var in pb.GetBasesInput
			ReadJsonInput(w, r, &in)
			filter_to_useful := in.FilterToUseful
			include_market_goods := in.IncludeMarketGoods

			var result []*configs_export.Base
			if filter_to_useful {
				result = configs_export.FitlerToUsefulOres(api.app_data.Configs.MiningOperations)
			} else {
				result = api.app_data.Configs.MiningOperations
			}
			result = services.FilterNicknames(in.FilterNicknames, result)

			var output []*Base
			for _, item := range result {
				answer := &Base{
					Base: item,
				}
				if include_market_goods {
					answer.MarketGoods = services.FilterMarketGoodCategory(in.FilterMarketGoodCategory, item.MarketGoodsPerNick)
				}
				output = append(output, answer)
			}
			apiutils.ReturnJson(&w, output)
		},
	}
}
