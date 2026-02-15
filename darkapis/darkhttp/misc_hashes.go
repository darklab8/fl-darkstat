package darkhttp

import (
	"net/http"

	"github.com/darklab8/fl-darkstat/darkapis/darkgrpc_deprecated"
	"github.com/darklab8/fl-darkstat/darkapis/darkhttp/apiutils"
	"github.com/darklab8/fl-darkstat/darkcore/web"
	"github.com/darklab8/fl-darkstat/darkcore/web/registry"
)

// ShowAccount godoc
// @Summary      Hashes
// @Tags         misc
// @Accept       json
// @Produce      json
// @Success      200  {object}  	darkgrpc_deprecated.Hashes
// @Router       /api/hashes [post]
func GetHashes(webapp *web.Web, api *Api) *registry.Endpoint {
	return &registry.Endpoint{
		Url: "" + ApiRoute + "/hashes",
		Handler: func(w http.ResponseWriter, r *http.Request) {
			if webapp.AppDataMutex != nil {
				webapp.AppDataMutex.RLock()
				defer webapp.AppDataMutex.RUnlock()
			}
			hashes := darkgrpc_deprecated.GetHashesData(api.app_data)
			apiutils.ReturnJson(&w, darkgrpc_deprecated.Hashes{HashesByNick: hashes})
		},
	}
}
func (c *HttpClient) GetHashes() (darkgrpc_deprecated.Hashes, error) {
	return make_request[EmptyInput, darkgrpc_deprecated.Hashes](c, ""+ApiRoute+"/hashes", EmptyInput{})
}
