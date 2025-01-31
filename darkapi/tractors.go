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
// @Summary      Getting list of tractors
// @Tags         tractors
// @Accept       json
// @Produce      json
// @Success      200  {array}  	configs_export.Tractor
// @Router       /api/tractors [get]
func GetTractors(webapp *web.Web, api *Api) *registry.Endpoint {
	return &registry.Endpoint{
		Url: "GET " + ApiRoute + "/tractors",
		Handler: func(w http.ResponseWriter, r *http.Request) {
			if webapp.AppDataMutex != nil {
				webapp.AppDataMutex.Lock()
				defer webapp.AppDataMutex.Unlock()
			}
			ReturnJson(&w, api.app_data.Configs.Tractors)
		},
	}
}

// ShowAccount godoc
// @Summary      Getting list of Tractor Market Goods
// @Tags         tractors
// @Accept       json
// @Produce      json
// @Param request body []string true "Array of ship nicknames as input, for example [ai_bomber]" example("ai_bomber")
// @Success      200  {array}  	MarketGoodResp
// @Router       /api/tractors/market_goods [post]
func PostTractorMarketGoods(webapp *web.Web, api *Api) *registry.Endpoint {
	return &registry.Endpoint{
		Url: "POST " + ApiRoute + "/tractors/market_goods",
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

			items_by_nick := make(map[string]*configs_export.Tractor)
			for _, item := range api.app_data.Configs.Tractors {
				items_by_nick[string(item.Nickname)] = item
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
