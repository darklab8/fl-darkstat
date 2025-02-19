package darkhttp

import (
	"net/http"

	"github.com/darklab8/fl-darkstat/darkapis/darkgrpc"
	pb "github.com/darklab8/fl-darkstat/darkapis/darkgrpc/statproto"
	"github.com/darklab8/fl-darkstat/darkapis/darkhttp/apiutils"
	"github.com/darklab8/fl-darkstat/darkcore/web"
	"github.com/darklab8/fl-darkstat/darkcore/web/registry"
	"github.com/darklab8/fl-darkstat/darkstat/configs_export"
)

type CounterMeasure struct {
	*configs_export.CounterMeasure
	MarketGoods []*configs_export.MarketGood        `json:"market_goods"`
	TechCompat  *configs_export.DiscoveryTechCompat `json:"tech_compat"`
}

// ShowAccount godoc
// @Summary      Getting list of CounterMeasure
// @Tags         equipment
// @Accept       json
// @Produce      json
// @Success      200  {array}  	darkhttp.CounterMeasure
// @Router       /api/counter_measures [post]
// @Param request body pb.GetEquipmentInput true "input variables"
// @Description  include_market_goods: "insert 'true' if wish to include market goods under 'market goods' key or not. Such data can add a lot of extra weight"
// @Description  include_tech_compat: insert 'true' if wish to include tech compatibility data. can be adding a lot of extra weight
// @Description  filter_to_useful: Apply filtering same as darkstat does by default for its tab. Usually means showing only items that can be bought/crafted/or found
// @Description  filter_nicknames: filters by item nicknames
func GetCMs(webapp *web.Web, api *Api) *registry.Endpoint {
	return &registry.Endpoint{
		Url: "GET " + ApiRoute + "/counter_measures",
		Handler: func(w http.ResponseWriter, r *http.Request) {
			if webapp.AppDataMutex != nil {
				webapp.AppDataMutex.Lock()
				defer webapp.AppDataMutex.Unlock()
			}

			var in *pb.GetEquipmentInput
			in, err := GetEquipmentInput(w, r)
			if err != nil {
				return
			}

			var result []configs_export.CounterMeasure
			if in.FilterToUseful {
				result = api.app_data.Configs.FilterToUsefulCounterMeasures(api.app_data.Configs.CMs)
			} else {
				result = api.app_data.Configs.CMs
			}
			result = darkgrpc.FilterNicknames(in.FilterNicknames, result)

			var output []*CounterMeasure
			for _, item := range result {
				answer := &CounterMeasure{
					CounterMeasure: &item,
				}
				if in.IncludeMarketGoods {
					for _, good := range item.Bases {
						answer.MarketGoods = append(answer.MarketGoods, good)
					}
				}
				if in.IncludeTechCompat {
					answer.TechCompat = item.DiscoveryTechCompat
				}
				output = append(output, answer)
			}

			apiutils.ReturnJson(&w, output)
		}}
}
