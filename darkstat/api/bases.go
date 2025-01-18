package api

import (
	"net/http"

	"github.com/darklab8/fl-darkcore/darkcore/web"
	"github.com/darklab8/fl-darkcore/darkcore/web/registry"
)

// ShowAccount godoc
// @Summary      Getting list of NPC Bases
// @Tags         bases
// @Accept       json
// @Produce      json
// @Success      200  {array}  	configs_export.Base
// @Router       /api/npc_bases [get]
func NewEndpointBases(webapp *web.Web, api *Api) *registry.Endpoint {
	return &registry.Endpoint{
		Url: "GET " + ApiRoute + "/npc_bases",
		Handler: func(w http.ResponseWriter, r *http.Request) {
			if webapp.AppDataMutex != nil {
				webapp.AppDataMutex.Lock()
				defer webapp.AppDataMutex.Unlock()
			}
			ReturnJson(&w, api.app_data.Configs.Bases)
		},
	}
}

// ShowAccount godoc
// @Summary      Getting list of Mining Operations
// @Tags         bases
// @Accept       json
// @Produce      json
// @Success      200  {array}  	configs_export.Base
// @Router       /api/mining_operations [get]
func NewEndpointOreFields(webapp *web.Web, api *Api) *registry.Endpoint {
	return &registry.Endpoint{
		Url: "GET " + ApiRoute + "/mining_operations",
		Handler: func(w http.ResponseWriter, r *http.Request) {
			if webapp.AppDataMutex != nil {
				webapp.AppDataMutex.Lock()
				defer webapp.AppDataMutex.Unlock()
			}
			ReturnJson(&w, api.app_data.Configs.MiningOperations)
		},
	}
}
