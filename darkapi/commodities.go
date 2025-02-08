package darkapi

import (
	"net/http"

	"github.com/darklab8/fl-darkstat/darkcore/web"
	"github.com/darklab8/fl-darkstat/darkcore/web/registry"
	"github.com/darklab8/fl-darkstat/darkstat/configs_export"
)

type Commodity struct {
	*configs_export.Commodity
	MarketGoods []*configs_export.MarketGood `json:"market_goods"`
}

// ShowAccount godoc
// @Summary      Getting list of Commodities
// @Tags         equipment
// @Accept       json
// @Produce      json
// @Success      200  {array}  	darkapi.Commodity
// @Router       /api/commodities [get]
// @Param        filter_to_useful    query     string  false  "insert 'true' if wish to filter items only to useful, usually they are sold, or have goods, or craftable or findable in loot, or bases that are flight reachable from manhattan"
// @Param        include_market_goods    query     string  false  "insert 'true' if wish to include market goods under 'market goods' key or not. Such data can add a lot of extra weight"
func GetCommodities(webapp *web.Web, api *Api) *registry.Endpoint {
	return &registry.Endpoint{
		Url: "GET " + ApiRoute + "/commodities",
		// Handler: GetItemsT(webapp, api.app_data.Configs.Commodities, api.app_data.Configs.FilterToUsefulCommodities),
		Handler: func(w http.ResponseWriter, r *http.Request) {
			if webapp.AppDataMutex != nil {
				webapp.AppDataMutex.Lock()
				defer webapp.AppDataMutex.Unlock()
			}

			filter_to_useful := r.URL.Query().Get("filter_to_useful") == "true"
			include_market_goods := r.URL.Query().Get("include_market_goods") == "true"

			var result []*configs_export.Commodity
			if filter_to_useful {
				result = api.app_data.Configs.FilterToUsefulCommodities(api.app_data.Configs.Commodities)
			} else {
				result = api.app_data.Configs.Commodities
			}

			var output []*Commodity
			for _, item := range result {
				answer := &Commodity{
					Commodity: item,
				}
				if include_market_goods {
					for _, good := range item.Bases {
						answer.MarketGoods = append(answer.MarketGoods, good)
					}
				}
				output = append(output, answer)
			}

			ReturnJson(&w, output)
		},
	}
}
