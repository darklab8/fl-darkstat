package darkhttp

import (
	"net/http"

	pb "github.com/darklab8/fl-darkstat/darkapis/darkgrpc/statproto"
	"github.com/darklab8/fl-darkstat/darkapis/darkhttp/apiutils"
	"github.com/darklab8/fl-darkstat/darkcore/web"
	"github.com/darklab8/fl-darkstat/darkcore/web/registry"
	"github.com/darklab8/fl-darkstat/darkstat/configs_export"
)

// ShowAccount godoc
// @Summary      Getting list of Factions
// @Tags         factions
// @Accept       json
// @Produce      json
// @Success      200  {array}  	configs_export.Faction
// @Router       /api/factions [post]
// @Param request body pb.GetFactionsInput true "input variables"
// @Description  filter_to_useful: Apply filtering same as darkstat does by default for its tab. Usually means showing only items that can be bought/crafted/or found
func GetFactions(webapp *web.Web, api *Api) *registry.Endpoint {
	return &registry.Endpoint{
		Url: "" + ApiRoute + "/factions",
		Handler: func(w http.ResponseWriter, r *http.Request) {
			if webapp.AppDataMutex != nil {
				webapp.AppDataMutex.RLock()
				defer webapp.AppDataMutex.RUnlock()
			}

			var in *pb.GetFactionsInput = &pb.GetFactionsInput{}
			if err := ReadJsonInput(w, r, &in); err != nil && r.Method == "POST" {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}
			if r.URL.Query().Get("filter_to_useful") == "true" {
				in.FilterToUseful = true
			}

			var input []configs_export.Faction
			if in.FilterToUseful {
				input = configs_export.FilterToUsefulFactions(api.app_data.Configs.Factions)
			} else {
				input = api.app_data.Configs.Factions
			}

			var result []configs_export.Faction
			for _, item := range input {
				copied := item
				if !in.IncludeBribes {
					copied.Bribes = []configs_export.Bribe{}
				}
				if !in.IncludeReputations {
					copied.Reputations = []configs_export.Reputation{}
				}
				result = append(result, copied)
			}

			apiutils.ReturnJson(&w, result)
		},
	}
}
