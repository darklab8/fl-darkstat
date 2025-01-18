package api

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/darklab8/fl-configs/configs/cfgtype"
	"github.com/darklab8/fl-configs/configs/configs_export/trades"
	"github.com/darklab8/fl-darkcore/darkcore/web"
	"github.com/darklab8/fl-darkcore/darkcore/web/registry"
	"github.com/darklab8/fl-darkstat/darkstat/settings/logus"
	"github.com/darklab8/go-utils/utils/ptr"
)

type GraphPathReq struct {
	From string `json:"from" example:"li01_01_base"` // Write NPC base nickname, or PoB nickname (Name in base64 encoding) or Ore field name
	To   string `json:"to" example:"br01_01_base"`   // Write NPC base nickname, or PoB nickname (Name in base64 encoding) or Ore field name
}

type GraphPathTime struct {
	Transport *cfgtype.SecondsI `json:"transport"` // time in seconds
	Frigate   *cfgtype.SecondsI `json:"frigate"`   // time in seconds
	Freighter *cfgtype.SecondsI `json:"freighter"` // time in seconds
}

type GraphPathsResp struct {
	Route GraphPathReq   `json:"route"` // writes requested input
	Time  *GraphPathTime `json:"time,omitempty"`
	Error *string        `json:"error,omitempty"` // writes error if requesting not existing nicknames in from/to fields
}

// ShowAccount godoc
// @Summary      List of time measurements between two NPC bases/PoBs/Ore fields and etc.
// @Description  You query by nicknames of objects from which base/pob/ore fields to which one
// @Description  You receive result how many seconds it takes to reach destination for Transport, Frigate and Freighter
// @Description  If destination is not reachable, you get time equal to Maximum of int32 = 9223372036854775807
// @Tags         graph
// @Accept       json
// @Produce      json
// @Param request body []api.GraphPathReq true "Request body"
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

			var input_routes []GraphPathReq
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

			var output_routes []GraphPathsResp

			for _, route := range input_routes {
				result := GraphPathsResp{Route: route}

				var transport_time, frigate_time, freighter_time int
				var err error
				transport_time, _ = trades.GetTimeMs(
					api.app_data.Configs.Transport.Graph,
					api.app_data.Configs.Transport.Time,
					route.From,
					route.To)
				frigate_time, _ = trades.GetTimeMs(
					api.app_data.Configs.Frigate.Graph,
					api.app_data.Configs.Frigate.Time,
					route.From,
					route.To)
				freighter_time, err = trades.GetTimeMs(
					api.app_data.Configs.Freighter.Graph,
					api.app_data.Configs.Freighter.Time,
					route.From,
					route.To)

				if err != nil {
					result.Error = ptr.Ptr(err.Error())
				} else {
					result.Time = &GraphPathTime{
						Transport: ptr.Ptr(transport_time / int(trades.PrecisionMultipiler)),
						Frigate:   ptr.Ptr(frigate_time / int(trades.PrecisionMultipiler)),
						Freighter: ptr.Ptr(freighter_time / int(trades.PrecisionMultipiler)),
					}
				}

				output_routes = append(output_routes, result)
			}

			data, err := json.Marshal(output_routes)
			logus.Log.CheckPanic(err, "should be marshable")
			fmt.Fprint(resp, string(data))
		},
	}
}
