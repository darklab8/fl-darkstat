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
// @Param request body pb.GetEquipmentInput true "input variables, description in Models of api 2.0"
func GetAmmos(webapp *web.Web, api *Api) *registry.Endpoint {
	return &registry.Endpoint{
		Url: "" + ApiRoute + "/ammos",
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

			var result []configs_export.Ammo
			if in.FilterToUseful {
				result = api.app_data.Configs.FilterToUsefulAmmo(api.app_data.Configs.Ammos)
			} else {
				result = api.app_data.Configs.Ammos
			}
			result = darkgrpc.FilterNicknames(in.FilterNicknames, result)

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
