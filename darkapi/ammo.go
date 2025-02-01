package darkapi

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/darklab8/fl-darkstat/darkcore/web"
	"github.com/darklab8/fl-darkstat/darkcore/web/registry"
	"github.com/darklab8/fl-darkstat/darkstat/configs_export"
	"github.com/darklab8/fl-darkstat/darkstat/settings/logus"
	"github.com/darklab8/go-utils/utils/ptr"
)

// ShowAccount godoc
// @Summary      Getting list of Ammos
// @Tags         ammos
// @Accept       json
// @Produce      json
// @Success      200  {array}  	configs_export.Ammo
// @Router       /api/ammos [get]
func GetAmmos(webapp *web.Web, api *Api) *registry.Endpoint {
	return &registry.Endpoint{
		Url: "GET " + ApiRoute + "/ammos",
		Handler: func(w http.ResponseWriter, r *http.Request) {
			if webapp.AppDataMutex != nil {
				webapp.AppDataMutex.Lock()
				defer webapp.AppDataMutex.Unlock()
			}
			ReturnJson(&w, api.app_data.Configs.Ammos)
		},
	}
}

// ShowAccount godoc
// @Summary      Getting list of Ammo Market Goods
// @Tags         ammos
// @Accept       json
// @Produce      json
// @Param request body []string true "Array of ammo nicknames as input, for example [ai_bomber]" example("ai_bomber")
// @Success      200  {array}  	MarketGoodResp
// @Router       /api/ammos/market_goods [post]
func PostAmmoMarketGoods(webapp *web.Web, api *Api) *registry.Endpoint {
	return &registry.Endpoint{
		Url: "POST " + ApiRoute + "/ammos/market_goods",
		Handler: func(w http.ResponseWriter, r *http.Request) {
			if webapp.AppDataMutex != nil {
				webapp.AppDataMutex.Lock()
				defer webapp.AppDataMutex.Unlock()
			}

			var nicknames []string
			body, err := io.ReadAll(r.Body)
			if logus.Log.CheckError(err, "failed to read body") {
				w.WriteHeader(http.StatusBadRequest)
				fmt.Fprintf(w, "err to ready body")
				return
			}
			json.Unmarshal(body, &nicknames)
			if len(nicknames) == 0 {
				w.WriteHeader(http.StatusBadRequest)
				fmt.Fprintf(w, "input at least some base nicknames into request body")
				return
			}

			var market_good_answers []*MarketGoodResp

			items_by_nick := make(map[string]*configs_export.Ammo)
			for _, item := range api.app_data.Configs.Ammos {
				items_by_nick[string(item.Nickname)] = &item
			}

			for _, input_nickname := range nicknames {
				answer := &MarketGoodResp{Nickname: string(input_nickname)}
				if ship, ok := items_by_nick[input_nickname]; ok {
					for _, good := range ship.Bases {
						answer.MarketGoods = append(answer.MarketGoods, good)
					}
				} else {
					answer.Error = ptr.Ptr("not existing ship nickname")
				}
				market_good_answers = append(market_good_answers, answer)

			}
			ReturnJson(&w, market_good_answers)
		},
	}
}

// ShowAccount godoc
// @Summary      Getting list of Ammos Tech compats
// @Tags         ammos
// @Accept       json
// @Produce      json
// @Param request body []string true "Array of ammo nicknames as input
// @Success      200  {array}  	TechCompatResp
// @Router       /api/ammos/tech_compats [post]
func PostAmmoTechcompatibilities(webapp *web.Web, api *Api) *registry.Endpoint {
	return &registry.Endpoint{
		Url: "POST " + ApiRoute + "/ammos/tech_compats",
		Handler: func(w http.ResponseWriter, r *http.Request) {
			if webapp.AppDataMutex != nil {
				webapp.AppDataMutex.Lock()
				defer webapp.AppDataMutex.Unlock()
			}

			var nicknames []string
			body, err := io.ReadAll(r.Body)
			if logus.Log.CheckError(err, "failed to read body") {
				w.WriteHeader(http.StatusBadRequest)
				fmt.Fprintf(w, "err to ready body")
				return
			}
			json.Unmarshal(body, &nicknames)
			if len(nicknames) == 0 {
				w.WriteHeader(http.StatusBadRequest)
				fmt.Fprintf(w, "input at least some base nicknames into request body")
				return
			}

			var market_good_answers []*TechCompatResp

			items_by_nicknames := make(map[string]*configs_export.Ammo)
			for _, item := range api.app_data.Configs.Ammos {
				items_by_nicknames[string(item.Nickname)] = &item
			}

			for _, input_nickname := range nicknames {
				answer := &TechCompatResp{Nickname: string(input_nickname)}
				if ship, ok := items_by_nicknames[input_nickname]; ok {
					answer.TechCompat = ship.DiscoveryTechCompat
				} else {
					answer.Error = ptr.Ptr("not existing ship nickname")
				}
				market_good_answers = append(market_good_answers, answer)

			}
			ReturnJson(&w, market_good_answers)
		},
	}
}
