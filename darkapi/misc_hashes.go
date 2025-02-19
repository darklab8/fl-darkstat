package darkapi

import (
	"net/http"

	"github.com/darklab8/fl-darkstat/darkapi/apiutils"
	"github.com/darklab8/fl-darkstat/darkcore/web"
	"github.com/darklab8/fl-darkstat/darkcore/web/registry"
	"github.com/darklab8/fl-darkstat/public/services"
)

// ShowAccount godoc
// @Summary      Hashes
// @Tags         misc
// @Accept       json
// @Produce      json
// @Success      200  {object}  	services.Hashes
// @Router       /api/hashes [get]
func GetHashes(webapp *web.Web, api *Api) *registry.Endpoint {
	return &registry.Endpoint{
		Url: "GET " + ApiRoute + "/hashes",
		Handler: func(w http.ResponseWriter, r *http.Request) {
			if webapp.AppDataMutex != nil {
				webapp.AppDataMutex.Lock()
				defer webapp.AppDataMutex.Unlock()
			}
			hashes := services.GetHashesData(api.app_data)
			apiutils.ReturnJson(&w, services.Hashes{HashesByNick: hashes})
		},
	}
}
