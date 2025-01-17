package api

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/darklab8/fl-darkcore/darkcore/web"
	"github.com/darklab8/fl-darkcore/darkcore/web/registry"
	"github.com/darklab8/fl-darkstat/darkstat/settings/logus"
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

			data, err := json.Marshal(api.app_data.Configs.Bases)
			logus.Log.CheckPanic(err, "should be marshable")
			fmt.Fprint(w, string(data))
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

			data, err := json.Marshal(api.app_data.Configs.MiningOperations)
			logus.Log.CheckPanic(err, "should be marshable")
			fmt.Fprint(w, string(data))
		},
	}
}
