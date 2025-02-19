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

type Ammo struct {
	*configs_export.Ammo
	MarketGoods []*configs_export.MarketGood        `json:"market_goods"`
	TechCompat  *configs_export.DiscoveryTechCompat `json:"tech_compat"`
}

// ShowAccount godoc
// @Summary      Getting list of Ammos
// @Tags         equipment
// @Accept       json
// @Produce      json
// @Success      200  {array}  	darkapi.Ammo
// @Router       /api/ammos [get]
// @Param        filter_to_useful    query     string  false  "insert 'true' if wish to filter items only to useful, usually they are sold, or have goods, or craftable or findable in loot, or bases that are flight reachable from manhattan"
// @Param        include_market_goods    query     string  false  "insert 'true' if wish to include market goods under 'market goods' key or not. Such data can add a lot of extra weight"
// @Param        include_tech_compat    query     string  false  "insert 'true' if wish to include tech compat info too for the item. Such data can add a lot of extra weight"
func GetAmmos(webapp *web.Web, api *Api) *registry.Endpoint {
	return &registry.Endpoint{
		Url: "" + ApiRoute + "/ammos",
		Handler: func(w http.ResponseWriter, r *http.Request) {
			if webapp.AppDataMutex != nil {
				webapp.AppDataMutex.Lock()
				defer webapp.AppDataMutex.Unlock()
			}

			var in pb.GetEquipmentInput
			if err := ReadJsonInput(w, r, &in); err != nil && r.Method == "POST" {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}
			filter_to_useful := in.FilterToUseful
			include_market_goods := in.IncludeMarketGoods
			include_tech_compat := in.IncludeTechCompat
			if r.URL.Query().Get("filter_to_useful") == "true" {
				filter_to_useful = true
			}
			if r.URL.Query().Get("include_market_goods") == "true" {
				include_market_goods = true
			}
			if r.URL.Query().Get("include_tech_compat") == "true" {
				include_tech_compat = true
			}

			var result []configs_export.Ammo
			if filter_to_useful {
				result = api.app_data.Configs.FilterToUsefulAmmo(api.app_data.Configs.Ammos)
			} else {
				result = api.app_data.Configs.Ammos
			}
			result = services.FilterNicknames(in.FilterNicknames, result)

			var output []*Ammo
			for _, item := range result {
				answer := &Ammo{
					Ammo: &item,
				}
				if include_market_goods {
					for _, good := range services.FilterMarketGoodCategory(in.FilterMarketGoodCategory, item.Bases) {
						answer.MarketGoods = append(answer.MarketGoods, good)
					}
				}
				if include_tech_compat {
					answer.TechCompat = item.DiscoveryTechCompat
				}
				output = append(output, answer)
			}

			apiutils.ReturnJson(&w, output)
		},
	}
}
