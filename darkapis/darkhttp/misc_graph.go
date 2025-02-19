package darkhttp

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/darklab8/fl-darkstat/darkapis/darkhttp/apiutils"
	"github.com/darklab8/fl-darkstat/darkcore/web"
	"github.com/darklab8/fl-darkstat/darkcore/web/registry"
	"github.com/darklab8/fl-darkstat/darkstat/appdata"
	"github.com/darklab8/fl-darkstat/darkstat/settings/logus"
)

// ShowAccount godoc
// @Summary      List of time measurements between two NPC bases/PoBs/Ore fields and etc.
// @Description  You query by nicknames of objects from which base/pob/ore fields to which one
// @Description  You receive result how many seconds it takes to reach destination for Transport, Frigate and Freighter
// @Description  If destination is not reachable, you get time equal to Maximum of int32 = 9223372036854775807
// @Tags         misc
// @Accept       json
// @Produce      json
// @Param request body []appdata.GraphPathReq true "Request body"
// @Success      200  {array}  	appdata.GraphPathsResp
// @Router       /api/graph/paths [post]
func PostGraphPaths(webapp *web.Web, api *Api) *registry.Endpoint {
	return &registry.Endpoint{
		Url: "POST " + ApiRoute + "/graph/paths",
		Handler: func(resp http.ResponseWriter, r *http.Request) {
			if webapp.AppDataMutex != nil {
				webapp.AppDataMutex.Lock()
				defer webapp.AppDataMutex.Unlock()
			}

			var input_routes []appdata.GraphPathReq
			body, err := io.ReadAll(r.Body)
			if logus.Log.CheckError(err, "failed to read body") {
				resp.WriteHeader(http.StatusBadRequest)
				fmt.Fprintf(resp, "err to ready body")
				return
			}
			json.Unmarshal(body, &input_routes)

			if len(input_routes) == 0 {
				resp.WriteHeader(http.StatusBadRequest)
				fmt.Fprintf(resp, "input at least some routes into request body")
				return
			}
			output_routes := api.app_data.GetGraphPaths(input_routes)

			apiutils.ReturnJson(&resp, output_routes)
		},
	}
}
