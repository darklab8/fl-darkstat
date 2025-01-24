package api

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/darklab8/fl-configs/configs/cfgtype"
	"github.com/darklab8/fl-darkcore/darkcore/web"
	"github.com/darklab8/fl-darkcore/darkcore/web/registry"
	"github.com/darklab8/fl-darkstat/darkstat/configs_export"
	"github.com/darklab8/fl-darkstat/darkstat/settings/logus"
	"github.com/darklab8/go-utils/utils/ptr"
)

// ShowAccount godoc
// @Summary      Getting list of NPC Bases
// @Tags         bases
// @Accept       json
// @Produce      json
// @Success      200  {array}  	configs_export.Base
// @Router       /api/npc_bases [get]
func GetBases(webapp *web.Web, api *Api) *registry.Endpoint {
	return &registry.Endpoint{
		Url: "GET " + ApiRoute + "/npc_bases",
		Handler: func(w http.ResponseWriter, r *http.Request) {
			if webapp.AppDataMutex != nil {
				webapp.AppDataMutex.Lock()
				defer webapp.AppDataMutex.Unlock()
			}
			ReturnJson(&w, api.app_data.Configs.Bases)
		},
	}
}

type MarketGoodResp struct {
	MarketGoods []*configs_export.MarketGood `json:"market_goods"`
	Nickname    string                       `json:"nickname"`
	Error       *string                      `json:"error"`
}

// ShowAccount godoc
// @Summary      Getting list of NPC Bases Market Goods
// @Tags         bases
// @Accept       json
// @Produce      json
// @Param request body []cfgtype.BaseUniNick true "Array of npc base nicknames as input, for example [li01_01_base]" example("li01_01_base")
// @Success      200  {array}  	MarketGoodResp
// @Router       /api/npc_bases/market_goods [post]
func PostBaseMarketGoods(webapp *web.Web, api *Api) *registry.Endpoint {
	return &registry.Endpoint{
		Url: "POST " + ApiRoute + "/npc_bases/market_goods",
		Handler: func(w http.ResponseWriter, r *http.Request) {
			if webapp.AppDataMutex != nil {
				webapp.AppDataMutex.Lock()
				defer webapp.AppDataMutex.Unlock()
			}

			var base_nicknames []cfgtype.BaseUniNick
			body, err := io.ReadAll(r.Body)
			if logus.Log.CheckError(err, "failed to read body") {
				w.WriteHeader(http.StatusBadRequest)
				fmt.Fprintf(w, "err to ready body")
				return
			}
			json.Unmarshal(body, &base_nicknames)
			if len(base_nicknames) == 0 {
				w.WriteHeader(http.StatusBadRequest)
				fmt.Fprintf(w, "input at least some base nicknames into request body")
				return
			}

			var base_market_goods []*MarketGoodResp

			bases_by_nick := make(map[string]*configs_export.Base)
			for _, base := range api.app_data.Configs.Bases {
				bases_by_nick[string(base.Nickname)] = base
			}

			for _, base_nickname := range base_nicknames {
				answer := &MarketGoodResp{Nickname: string(base_nickname)}
				if base, ok := bases_by_nick[base_nickname.ToStr()]; ok {
					for _, good := range base.MarketGoodsPerNick {
						answer.MarketGoods = append(answer.MarketGoods, good)
					}
				} else {
					answer.Error = ptr.Ptr("not existing base")
				}
				base_market_goods = append(base_market_goods, answer)

			}
			ReturnJson(&w, base_market_goods)
		},
	}
}

// ShowAccount godoc
// @Summary      Getting list of Mining Operations
// @Tags         bases
// @Accept       json
// @Produce      json
// @Success      200  {array}  	configs_export.Base
// @Router       /api/mining_operations [get]
func GetOreFields(webapp *web.Web, api *Api) *registry.Endpoint {
	return &registry.Endpoint{
		Url: "GET " + ApiRoute + "/mining_operations",
		Handler: func(w http.ResponseWriter, r *http.Request) {
			if webapp.AppDataMutex != nil {
				webapp.AppDataMutex.Lock()
				defer webapp.AppDataMutex.Unlock()
			}
			ReturnJson(&w, api.app_data.Configs.MiningOperations)
		},
	}
}
