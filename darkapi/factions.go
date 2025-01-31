package darkapi

import (
	"net/http"

	"github.com/darklab8/fl-darkstat/darkcore/web"
	"github.com/darklab8/fl-darkstat/darkcore/web/registry"
)

// ShowAccount godoc
// @Summary      Getting list of Factions
// @Tags         factions
// @Accept       json
// @Produce      json
// @Success      200  {array}  	configs_export.Faction
// @Router       /api/factions [get]
func GetFactions(webapp *web.Web, api *Api) *registry.Endpoint {
	return &registry.Endpoint{
		Url: "GET " + ApiRoute + "/factions",
		Handler: func(w http.ResponseWriter, r *http.Request) {
			if webapp.AppDataMutex != nil {
				webapp.AppDataMutex.Lock()
				defer webapp.AppDataMutex.Unlock()
			}
			ReturnJson(&w, api.app_data.Configs.Factions)
		},
	}
}
