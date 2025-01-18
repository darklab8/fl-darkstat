package api

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/darklab8/fl-configs/configs/configs_mapped/freelancer_mapped/data_mapped/initialworld/flhash"
	"github.com/darklab8/fl-darkcore/darkcore/web"
	"github.com/darklab8/fl-darkcore/darkcore/web/registry"
	"github.com/darklab8/fl-darkstat/darkstat/settings/logus"
)

type Hashes struct {
	NicknameToHash map[string]flhash.HashCode `json:"nickname_to_hash"`
}

// ShowAccount godoc
// @Summary      Hashes
// @Tags         hashes
// @Accept       json
// @Produce      json
// @Success      200  {object}  	api.Hashes
// @Router       /api/hashes [get]
func NewEndpointHashes(webapp *web.Web, api *Api) *registry.Endpoint {
	return &registry.Endpoint{
		Url: "GET " + ApiRoute + "/hashes",
		Handler: func(w http.ResponseWriter, r *http.Request) {
			if webapp.AppDataMutex != nil {
				webapp.AppDataMutex.Lock()
				defer webapp.AppDataMutex.Unlock()
			}

			data, err := json.Marshal(Hashes{NicknameToHash: api.app_data.Configs.Hashes})
			logus.Log.CheckPanic(err, "should be marshable")
			fmt.Fprint(w, string(data))
		},
	}
}
