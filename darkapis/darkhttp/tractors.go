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

type Tractor struct {
	*configs_export.Tractor
	MarketGoods []*configs_export.MarketGood        `json:"market_goods"`
	TechCompat  *configs_export.DiscoveryTechCompat `json:"tech_compat"`
}

// ShowAccount godoc
// @Summary      Getting list of tractors
// @Tags         equipment
// @Accept       json
// @Produce      json
// @Success      200  {array}  	darkhttp.Tractor
// @Router       /api/tractors [post]
// @Param request body pb.GetTractorsInput true "input variables"
// @Description  include_market_goods: "insert 'true' if wish to include market goods under 'market goods' key or not. Such data can add a lot of extra weight"
// @Description  filter_to_useful: Apply filtering same as darkstat does by default for its tab. Usually means showing only items that can be bought/crafted/or found
// @Description  filter_nicknames: filters by item nicknames
func GetTractors(webapp *web.Web, api *Api) *registry.Endpoint {
	return &registry.Endpoint{
		Url: "" + ApiRoute + "/tractors",
		Handler: func(w http.ResponseWriter, r *http.Request) {
			if webapp.AppDataMutex != nil {
				webapp.AppDataMutex.Lock()
				defer webapp.AppDataMutex.Unlock()
			}

			var in *pb.GetTractorsInput = &pb.GetTractorsInput{}
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

			var result []*configs_export.Tractor
			if in.FilterToUseful {
				result = api.app_data.Configs.FilterToUsefulTractors(api.app_data.Configs.Tractors)
			} else {
				result = api.app_data.Configs.Tractors
			}
			result = darkgrpc.FilterNicknames(in.FilterNicknames, result)

			var output []*Tractor
			for _, item := range result {
				answer := &Tractor{
					Tractor: item,
				}
				if in.IncludeMarketGoods {
					for _, good := range item.Bases {
						answer.MarketGoods = append(answer.MarketGoods, good)
					}
				}
				output = append(output, answer)
			}

			apiutils.ReturnJson(&w, output)
		}}
}
