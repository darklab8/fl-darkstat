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

type Ammo struct {
	*configs_export.Ammo
	MarketGoods []*configs_export.MarketGood        `json:"market_goods"`
	TechCompat  *configs_export.DiscoveryTechCompat `json:"tech_compat"`
}

func GetEquipmentInput(w http.ResponseWriter, r *http.Request) (*pb.GetEquipmentInput, error) {
	var in *pb.GetEquipmentInput = &pb.GetEquipmentInput{}
	if err := ReadJsonInput(w, r, &in); err != nil && r.Method == "POST" {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return in, err
	}
	if r.URL.Query().Get("filter_to_useful") == "true" {
		in.FilterToUseful = true
	}
	if r.URL.Query().Get("include_market_goods") == "true" {
		in.IncludeMarketGoods = true
	}
	if r.URL.Query().Get("include_tech_compat") == "true" {
		in.IncludeTechCompat = true
	}
	return in, nil
}

// ShowAccount godoc
// @Summary      Getting list of Ammos
// @Tags         equipment
// @Accept       json
// @Produce      json
// @Success      200  {array}  	darkhttp.Ammo
// @Router       /api/ammos [post]
// @Param request body pb.GetEquipmentInput true "input variables"
// @Description  include_market_goods: "insert 'true' if wish to include market goods under 'market goods' key or not. Such data can add a lot of extra weight"
// @Description  include_tech_compat: insert 'true' if wish to include tech compatibility data. can be adding a lot of extra weight
// @Description  filter_to_useful: Apply filtering same as darkstat does by default for its tab. Usually means showing only items that can be bought/crafted/or found
// @Description  filter_nicknames: filters by item nicknames
func GetAmmos(webapp *web.Web, api *Api) *registry.Endpoint {
	return &registry.Endpoint{
		Url: "" + ApiRoute + "/ammos",
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

			var result []configs_export.Ammo
			if in.FilterToUseful {
				result = api.app_data.Configs.FilterToUsefulAmmo(api.app_data.Configs.Ammos)
			} else {
				result = api.app_data.Configs.Ammos
			}
			result = darkgrpc_deprecated.FilterNicknames(in.FilterNicknames, result)

			var output []*Ammo
			for _, item := range result {
				answer := &Ammo{
					Ammo: &item,
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
		},
	}
}
func (c *HttpClient) GetAmmos(input pb.GetCommoditiesInput) ([]*Ammo, error) {
	return make_request[pb.GetCommoditiesInput, []*Ammo](c, ""+ApiRoute+"/ammos", input)
}
