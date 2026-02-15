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
func GetGraphPaths(webapp *web.Web, api *Api) *registry.Endpoint {
	return &registry.Endpoint{
		Url: "POST " + ApiRoute + "/graph/paths",
		Handler: func(resp http.ResponseWriter, r *http.Request) {
			if webapp.AppDataMutex != nil {
				webapp.AppDataMutex.RLock()
				defer webapp.AppDataMutex.RUnlock()
			}

			var input_routes []appdata.GraphPathReq
			body, err := io.ReadAll(r.Body)
			if logus.Log.CheckError(err, "failed to read body") {
				resp.WriteHeader(http.StatusBadRequest)
				_, err = fmt.Fprintf(resp, "err to ready body")
				Log.CheckError(err, "fprintf post graph paths error")
				return
			}
			err = json.Unmarshal(body, &input_routes)
			Log.CheckWarn(err, "failed to unparmshal input in post graph paths")

			if len(input_routes) == 0 {
				resp.WriteHeader(http.StatusBadRequest)
				_, err = fmt.Fprintf(resp, "input at least some routes into request body")
				Log.CheckError(err, "fprintf post graph paths error 2")
				return
			}
			output_routes := api.app_data.GetGraphPaths(input_routes)

			apiutils.ReturnJson(&resp, output_routes)
		},
	}
}
func (c *HttpClient) GetGraphPaths(requests []appdata.GraphPathReq) ([]appdata.GraphPathsResp, error) {
	return make_request[[]appdata.GraphPathReq, []appdata.GraphPathsResp](c, ""+ApiRoute+"/graph/paths", requests)
}
