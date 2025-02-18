package appdata

import (
	"github.com/darklab8/fl-darkstat/configs/cfg"
	"github.com/darklab8/fl-darkstat/darkstat/configs_export/trades"
	"github.com/darklab8/go-utils/utils/ptr"
)

type GraphPathReq struct {
	From string `json:"from" example:"li01_01_base" validate:"required"` // Write NPC base nickname, or PoB nickname (Name in base64 encoding) or Ore field name
	To   string `json:"to" example:"br01_01_base" validate:"required"`   // Write NPC base nickname, or PoB nickname (Name in base64 encoding) or Ore field name
}

type GraphPathTime struct {
	Transport *cfg.SecondsI `json:"transport"` // time in seconds
	Frigate   *cfg.SecondsI `json:"frigate"`   // time in seconds
	Freighter *cfg.SecondsI `json:"freighter"` // time in seconds
}

type GraphPathsResp struct {
	Query GraphPathReq   `json:"route" validate:"required"` // writes requested input
	Time  *GraphPathTime `json:"time,omitempty"`
	Error *string        `json:"error,omitempty"` // writes error if requesting not existing nicknames in from/to fields
}

func (app_data *AppData) GetGraphPaths(input_routes []GraphPathReq) []GraphPathsResp {
	var output_routes []GraphPathsResp

	for _, route := range input_routes {
		result := GraphPathsResp{Query: route}

		var transport_time, frigate_time, freighter_time uint32
		var err error
		transport_time, _ = trades.GetTimeMs(
			app_data.Configs.Transport.Graph,
			app_data.Configs.Transport.Time,
			route.From,
			route.To)
		frigate_time, _ = trades.GetTimeMs(
			app_data.Configs.Frigate.Graph,
			app_data.Configs.Frigate.Time,
			route.From,
			route.To)
		freighter_time, err = trades.GetTimeMs(
			app_data.Configs.Freighter.Graph,
			app_data.Configs.Freighter.Time,
			route.From,
			route.To)

		if err != nil {
			result.Error = ptr.Ptr(err.Error())
		} else {
			result.Time = &GraphPathTime{
				Transport: ptr.Ptr(int(transport_time) / int(trades.PrecisionMultipiler)),
				Frigate:   ptr.Ptr(int(frigate_time) / int(trades.PrecisionMultipiler)),
				Freighter: ptr.Ptr(int(freighter_time) / int(trades.PrecisionMultipiler)),
			}
		}

		output_routes = append(output_routes, result)
	}
	return output_routes
}
