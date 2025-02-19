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

type Gun struct {
	*configs_export.Gun
	MarketGoods []*configs_export.MarketGood        `json:"market_goods"`
	TechCompat  *configs_export.DiscoveryTechCompat `json:"tech_compat"`
}

func GetGunsInput(w http.ResponseWriter, r *http.Request) (*pb.GetGunsInput, error) {
	var in *pb.GetGunsInput = &pb.GetGunsInput{}
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
	if r.URL.Query().Get("include_damage_bonuses") == "true" {
		in.IncludeDamageBonuses = true
	}
	return in, nil
}

// ShowAccount godoc
// @Summary      Getting list of Guns
// @Tags         equipment
// @Accept       json
// @Produce      json
// @Success      200  {array}  	darkhttp.Gun
// @Router       /api/guns [post]
// @Param request body pb.GetGunsInput true "input variables"
// @Description  include_market_goods: "insert 'true' if wish to include market goods under 'market goods' key or not. Such data can add a lot of extra weight"
// @Description  include_tech_compat: insert 'true' if wish to include tech compatibility data. can be adding a lot of extra weight
// @Description  filter_to_useful: Apply filtering same as darkstat does by default for its tab. Usually means showing only items that can be bought/crafted/or found
// @Description  filter_nicknames: filters by item nicknames
// @Description  include_damage_bonuses: insert 'true' if u wish added damage bonuses against specific shield types
func GetGuns(webapp *web.Web, api *Api) *registry.Endpoint {
	return &registry.Endpoint{
		Url:     "" + ApiRoute + "/guns",
		Handler: GunHandler(webapp, api, api.app_data.Configs.Guns)}
}

func GunHandler(webapp *web.Web, api *Api, guns []configs_export.Gun) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		if webapp.AppDataMutex != nil {
			webapp.AppDataMutex.Lock()
			defer webapp.AppDataMutex.Unlock()
		}

		var in *pb.GetGunsInput
		in, err := GetGunsInput(w, r)
		if err != nil {
			return
		}

		var result []configs_export.Gun
		if in.FilterToUseful {
			result = api.app_data.Configs.FilterToUsefulGun(guns)
		} else {
			result = guns
		}
		result = darkgrpc.FilterNicknames(in.FilterNicknames, result)

		var output []*Gun
		for _, item := range result {
			answer := &Gun{
				Gun: &item,
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
	}
}

// ShowAccount godoc
// @Summary      Getting list of Missiles
// @Tags         equipment
// @Accept       json
// @Produce      json
// @Success      200  {array}  	darkhttp.Gun
// @Router       /api/missiles [post]
// @Param request body pb.GetGunsInput true "input variables"
// @Description  include_market_goods: "insert 'true' if wish to include market goods under 'market goods' key or not. Such data can add a lot of extra weight"
// @Description  include_tech_compat: insert 'true' if wish to include tech compatibility data. can be adding a lot of extra weight
// @Description  filter_to_useful: Apply filtering same as darkstat does by default for its tab. Usually means showing only items that can be bought/crafted/or found
// @Description  filter_nicknames: filters by item nicknames
// @Description  include_damage_bonuses: insert 'true' if u wish added damage bonuses against specific shield types
func GetMissiles(webapp *web.Web, api *Api) *registry.Endpoint {
	return &registry.Endpoint{
		Url:     "" + ApiRoute + "/missiles",
		Handler: GunHandler(webapp, api, api.app_data.Configs.Missiles),
	}
}
