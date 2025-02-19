package darkhttp

import (
	"net/http"

	"github.com/darklab8/fl-darkstat/darkapis/darkgrpc"
	"github.com/darklab8/fl-darkstat/darkapis/darkhttp/apiutils"
	"github.com/darklab8/fl-darkstat/darkcore/web"
	"github.com/darklab8/fl-darkstat/darkcore/web/registry"
)

// ShowAccount godoc
// @Summary      Hashes
// @Tags         misc
// @Accept       json
// @Produce      json
// @Success      200  {object}  	darkgrpc.Hashes
// @Router       /api/hashes [get]
func GetHashes(webapp *web.Web, api *Api) *registry.Endpoint {
	return &registry.Endpoint{
		Url: "GET " + ApiRoute + "/hashes",
		Handler: func(w http.ResponseWriter, r *http.Request) {
			if webapp.AppDataMutex != nil {
				webapp.AppDataMutex.Lock()
				defer webapp.AppDataMutex.Unlock()
			}
			hashes := darkgrpc.GetHashesData(api.app_data)
			apiutils.ReturnJson(&w, darkgrpc.Hashes{HashesByNick: hashes})
		},
	}
}
