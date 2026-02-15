package darkhttp

import (
	"net/http"

	pb "github.com/darklab8/fl-darkstat/darkapis/darkgrpc_deprecated/statproto_deprecated"
	"github.com/darklab8/fl-darkstat/darkapis/darkhttp/apiutils"
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
// @Success      200  {array}  	darkhttp.Commodity
// @Router       /api/commodities [post]
// @Param request body pb.GetCommoditiesInput true "input variables"
// @Description  include_market_goods: "insert 'true' if wish to include market goods under 'market goods' key or not. Such data can add a lot of extra weight"
// @Description  filter_to_useful: Apply filtering same as darkstat does by default for its tab. Usually means showing only items that can be bought/crafted/or found
// @Description  filter_nicknames: filters by item nicknames
func GetCommodities(webapp *web.Web, api *Api) *registry.Endpoint {
	return &registry.Endpoint{
		Url: "" + ApiRoute + "/commodities",
		// Handler: GetItemsT(webapp, api.app_data.Configs.Commodities, api.app_data.Configs.FilterToUsefulCommodities),
		Handler: func(w http.ResponseWriter, r *http.Request) {
			if webapp.AppDataMutex != nil {
				webapp.AppDataMutex.RLock()
				defer webapp.AppDataMutex.RUnlock()
			}

			var in *pb.GetCommoditiesInput = &pb.GetCommoditiesInput{}
			if err := ReadJsonInput(w, r, &in); err != nil && r.Method == "POST" {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}
			if r.URL.Query().Get("filter_to_useful") == "true" {
				in.FilterToUseful = true
			}
			if r.URL.Query().Get("include_market_goods") == "true" {
				in.IncludeMarketGoods = true
			}

			var result []*configs_export.Commodity
			if in.FilterToUseful {
				result = api.app_data.Configs.FilterToUsefulCommodities(api.app_data.Configs.Commodities)
			} else {
				result = api.app_data.Configs.Commodities
			}

			var output []*Commodity
			for _, item := range result {
				answer := &Commodity{
					Commodity: item,
				}
				if in.IncludeMarketGoods {
					for _, good := range item.Bases {
						answer.MarketGoods = append(answer.MarketGoods, good)
					}
				}
				output = append(output, answer)
			}

			apiutils.ReturnJson(&w, output)
		},
	}
}
func (c *HttpClient) GetCommodities(input pb.GetCommoditiesInput) ([]*Commodity, error) {
	return make_request[pb.GetCommoditiesInput, []*Commodity](c, ""+ApiRoute+"/commodities", input)
}
