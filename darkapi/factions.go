package darkapi

import (
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
		Url:     "GET " + ApiRoute + "/factions",
		Handler: GetItemsT(webapp, api.app_data.Configs.Factions),
	}
}
