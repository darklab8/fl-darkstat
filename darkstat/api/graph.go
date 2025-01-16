package api

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/darklab8/fl-configs/configs/configs_export/trades"
	"github.com/darklab8/fl-darkcore/darkcore/web"
	"github.com/darklab8/fl-darkcore/darkcore/web/registry"
	"github.com/darklab8/fl-darkstat/darkstat/settings/logus"
	"github.com/darklab8/go-utils/utils/ptr"
)

type GraphPathBody struct {
	Routes []GraphPathReq `json:"routes"`
}

type GraphPathReq struct {
	From string `json:"from" example:"li01_01_base"`
	To   string `json:"to" example:"br01_01_base"`
}

type GraphPathTime struct {
	Transport *int `json:"transport"`
	Frigate   *int `json:"frigate"`
	Freighter *int `json:"freighter"`
}

type GraphPathsResp struct {
	Route GraphPathReq  `json:"route"`
	Time  GraphPathTime `json:"time"`
}

// ShowAccount godoc
// @Summary      List of time measurements between two NPC bases/PoBs/Ore fields and etc.
// @Description  You query by nicknames of objects from which base/pob/ore fields to which one
// @Description  You receive result how many seconds it takes to reach destination for Transport, Frigate and Freighter
// @Description  If destination is not reachable, you get time equal to Maximum of int32 = 9223372036854775807
// @Tags         graph
// @Accept       json
// @Produce      json
// @Param request body api.GraphPathBody true "Request body"
// @Success      200  {array}  	api.GraphPathsResp
// @Router       /api/graph/paths [post]
func PostGraphPaths(webapp *web.Web, api *Api) *registry.Endpoint {
	return &registry.Endpoint{
		Url: "POST " + ApiRoute + "/graph/paths",
		Handler: func(resp http.ResponseWriter, r *http.Request) {
			if webapp.AppDataMutex != nil {
				webapp.AppDataMutex.Lock()
				defer webapp.AppDataMutex.Unlock()
			}

			var input_routes GraphPathBody
			body, err := io.ReadAll(r.Body)
			if logus.Log.CheckError(err, "failed to read body") {
				resp.WriteHeader(http.StatusBadRequest)
				fmt.Fprintf(resp, "err to ready body")
				return
			}
			json.Unmarshal(body, &input_routes)

			var output_routes []GraphPathsResp

			for _, route := range input_routes.Routes {
				result := GraphPathsResp{Route: route}

				result.Time.Transport = ptr.Ptr(trades.GetDist(
					api.app_data.Configs.Transport.Graph,
					api.app_data.Configs.Transport.Time,
					route.From,
					route.To))
				result.Time.Frigate = ptr.Ptr(trades.GetDist(
					api.app_data.Configs.Frigate.Graph,
					api.app_data.Configs.Frigate.Time,
					route.From,
					route.To))
				result.Time.Freighter = ptr.Ptr(trades.GetDist(
					api.app_data.Configs.Freighter.Graph,
					api.app_data.Configs.Freighter.Time,
					route.From,
					route.To))

				output_routes = append(output_routes, result)
			}

			data, err := json.Marshal(output_routes)
			logus.Log.CheckPanic(err, "should be marshable")
			fmt.Fprint(resp, string(data))
		},
	}
}
