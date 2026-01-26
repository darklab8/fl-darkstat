package darkhttp

import (
	"net/http"

	"github.com/darklab8/fl-darkstat/darkapis/darkhttp/apiutils"
	"github.com/darklab8/fl-darkstat/darkcore/web"
	"github.com/darklab8/fl-darkstat/darkcore/web/registry"
	"github.com/darklab8/fl-darkstat/darkmap/front/export_front"
)

// ShowAccount godoc
// @Summary      Getting list of Systems
// @Tags         equipment
// @Accept       json
// @Produce      json
// @Success      200  {array}  	export_front.System
// @Router       /api/systems [post]
func GetSystems(webapp *web.Web, api *Api) *registry.Endpoint {
	return &registry.Endpoint{
		Url: "" + ApiRoute + "/systems",
		Handler: func(w http.ResponseWriter, r *http.Request) {
			if webapp.AppDataMutex != nil {
				webapp.AppDataMutex.RLock()
				defer webapp.AppDataMutex.RUnlock()
			}

			var systems []export_front.System

			systems = export_front.ExportSystems(api.app_data.Configs.Mapped)

			apiutils.ReturnJson(&w, systems)
		}}
}
