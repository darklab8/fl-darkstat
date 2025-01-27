package api

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
// @Summary      Getting list of Ships
// @Tags         ships
// @Accept       json
// @Produce      json
// @Success      200  {array}  	configs_export.Ship
// @Router       /api/ships [get]
func GetShips(webapp *web.Web, api *Api) *registry.Endpoint {
	return &registry.Endpoint{
		Url: "GET " + ApiRoute + "/ships",
		Handler: func(w http.ResponseWriter, r *http.Request) {
			if webapp.AppDataMutex != nil {
				webapp.AppDataMutex.Lock()
				defer webapp.AppDataMutex.Unlock()
			}
			ReturnJson(&w, api.app_data.Configs.Ships)
		},
	}
}

// ShowAccount godoc
// @Summary      Getting list of Ship Market Goods
// @Tags         ships
// @Accept       json
// @Produce      json
// @Param request body []string true "Array of ship nicknames as input, for example [ai_bomber]" example("ai_bomber")
// @Success      200  {array}  	MarketGoodResp
// @Router       /api/ships/market_goods [post]
func PostShipMarketGoods(webapp *web.Web, api *Api) *registry.Endpoint {
	return &registry.Endpoint{
		Url: "POST " + ApiRoute + "/ships/market_goods",
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

			ships_by_nick := make(map[string]*configs_export.Ship)
			for _, item := range api.app_data.Configs.Ships {
				ships_by_nick[string(item.Nickname)] = &item
			}

			for _, input_nickname := range nicknames {
				answer := &MarketGoodResp{Nickname: string(input_nickname)}
				if ship, ok := ships_by_nick[input_nickname]; ok {
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

type TechCompatResp struct {
	TechCompat *configs_export.DiscoveryTechCompat `json:"tech_compat"`
	Nickname   string                              `json:"nickname"`
	Error      *string                             `json:"error"`
}

// ShowAccount godoc
// @Summary      Getting list of Ship Tech compats
// @Tags         ships
// @Accept       json
// @Produce      json
// @Param request body []string true "Array of ship nicknames as input, for example [ai_bomber]" example("ai_bomber")
// @Success      200  {array}  	TechCompatResp
// @Router       /api/ships/tech_compats [post]
func PostShipTechcompatibilities(webapp *web.Web, api *Api) *registry.Endpoint {
	return &registry.Endpoint{
		Url: "POST " + ApiRoute + "/ships/tech_compats",
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

			ships_by_nick := make(map[string]*configs_export.Ship)
			for _, item := range api.app_data.Configs.Ships {
				ships_by_nick[string(item.Nickname)] = &item
			}

			for _, input_nickname := range nicknames {
				answer := &TechCompatResp{Nickname: string(input_nickname)}
				if ship, ok := ships_by_nick[input_nickname]; ok {
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
