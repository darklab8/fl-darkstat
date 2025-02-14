package darkapi

import (
	"net/http"

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
// @Router       /api/npc_bases [get]
// @Param        filter_to_useful    query     string  false  "insert 'true' if wish to filter items only to useful, usually they are sold, or have goods, or craftable or findable in loot, or bases that are flight reachable from manhattan"
// @Param        include_market_goods    query     string  false  "insert 'true' if wish to include market goods under 'market goods' key or not. Such data can add a lot of extra weight"
func GetBases(webapp *web.Web, api *Api) *registry.Endpoint {
	return &registry.Endpoint{
		Url: "GET " + ApiRoute + "/npc_bases",
		Handler: func(w http.ResponseWriter, r *http.Request) {
			if webapp.AppDataMutex != nil {
				webapp.AppDataMutex.Lock()
				defer webapp.AppDataMutex.Unlock()
			}
			filter_to_useful := r.URL.Query().Get("filter_to_useful") == "true"
			include_market_goods := r.URL.Query().Get("include_market_goods") == "true"

			var result []*configs_export.Base
			if filter_to_useful {
				result = configs_export.FilterToUserfulBases(api.app_data.Configs.Bases)
			} else {
				result = api.app_data.Configs.Bases
			}

			var output []*Base
			for _, item := range result {
				answer := &Base{
					Base: item,
				}
				if include_market_goods {
					for _, good := range item.MarketGoodsPerNick {
						answer.MarketGoods = append(answer.MarketGoods, good)
					}
				}
				output = append(output, answer)
			}
			ReturnJson(&w, output)
		},
	}

}

// ShowAccount godoc
// @Summary      Getting list of Mining Operations
// @Tags         bases
// @Accept       json
// @Produce      json
// @Success      200  {array}  	darkapi.Base
// @Router       /api/mining_operations [get]
// @Param        filter_to_useful    query     string  false  "filter items only to useful, usually they are sold, or have goods, or craftable or findable in loot, or bases that are flight reachable from manhattan"
// @Param        include_market_goods    query     string  false  "include market goods under 'market goods' key or not. Such data can add a lot of extra weight"
func GetOreFields(webapp *web.Web, api *Api) *registry.Endpoint {
	return &registry.Endpoint{
		Url: "GET " + ApiRoute + "/mining_operations",
		Handler: func(w http.ResponseWriter, r *http.Request) {
			if webapp.AppDataMutex != nil {
				webapp.AppDataMutex.Lock()
				defer webapp.AppDataMutex.Unlock()
			}
			filter_to_useful := r.URL.Query().Get("filter_to_useful") == "true"
			include_market_goods := r.URL.Query().Get("include_market_goods") == "true"

			var result []*configs_export.Base
			if filter_to_useful {
				result = configs_export.FitlerToUsefulOres(api.app_data.Configs.MiningOperations)
			} else {
				result = api.app_data.Configs.MiningOperations
			}

			var output []*Base
			for _, item := range result {
				answer := &Base{
					Base: item,
				}
				if include_market_goods {
					for _, good := range item.MarketGoodsPerNick {
						answer.MarketGoods = append(answer.MarketGoods, good)
					}
				}
				output = append(output, answer)
			}
			ReturnJson(&w, output)
		},
	}
}
