package darkhttp

import (
	"net/http"

	"github.com/darklab8/fl-darkstat/darkapis/darkgrpc_deprecated"
	pb "github.com/darklab8/fl-darkstat/darkapis/darkgrpc_deprecated/statproto_deprecated"
	"github.com/darklab8/fl-darkstat/darkapis/darkhttp/apiutils"
	"github.com/darklab8/fl-darkstat/darkcore/web"
	"github.com/darklab8/fl-darkstat/darkcore/web/registry"
	"github.com/darklab8/fl-darkstat/darkstat/configs_export"
)

type Mine struct {
	*configs_export.Mine
	MarketGoods []*configs_export.MarketGood        `json:"market_goods"`
	TechCompat  *configs_export.DiscoveryTechCompat `json:"tech_compat"`
}

// ShowAccount godoc
// @Summary      Getting list of Mines
// @Tags         equipment
// @Accept       json
// @Produce      json
// @Success      200  {array}  	darkhttp.Mine
// @Router       /api/mines [post]
// @Param request body pb.GetEquipmentInput true "input variables"
// @Description  include_market_goods: "insert 'true' if wish to include market goods under 'market goods' key or not. Such data can add a lot of extra weight"
// @Description  include_tech_compat: insert 'true' if wish to include tech compatibility data. can be adding a lot of extra weight
// @Description  filter_to_useful: Apply filtering same as darkstat does by default for its tab. Usually means showing only items that can be bought/crafted/or found
// @Description  filter_nicknames: filters by item nicknames
func GetMines(webapp *web.Web, api *Api) *registry.Endpoint {
	return &registry.Endpoint{
		Url: "" + ApiRoute + "/mines",
		Handler: func(w http.ResponseWriter, r *http.Request) {
			if webapp.AppDataMutex != nil {
				webapp.AppDataMutex.RLock()
				defer webapp.AppDataMutex.RUnlock()
			}

			var in *pb.GetEquipmentInput
			in, err := GetEquipmentInput(w, r)
			if err != nil {
				return
			}

			var result []configs_export.Mine
			if in.FilterToUseful {
				result = api.app_data.Configs.FilterToUsefulMines(api.app_data.Configs.Mines)
			} else {
				result = api.app_data.Configs.Mines
			}
			result = darkgrpc_deprecated.FilterNicknames(in.FilterNicknames, result)

			var output []*Mine
			for _, item := range result {
				answer := &Mine{
					Mine: &item,
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
