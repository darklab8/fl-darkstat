package darkapi

import (
	"github.com/darklab8/fl-darkstat/darkcore/web"
	"github.com/darklab8/fl-darkstat/darkcore/web/registry"
	"github.com/darklab8/fl-darkstat/darkstat/configs_export"
)

type CmWithInfocard struct {
	configs_export.CounterMeasure
	Infocard configs_export.Infocard
}

// ShowAccount godoc
// @Summary      Getting list of CounterMeasure
// @Tags         counter_measures
// @Accept       json
// @Produce      json
// @Success      200  {array}  	CmWithInfocard
// @Router       /api/counter_measures [get]
// @Param        filter_to_useful    query     string  false  "filter items only to useful, usually they are sold, or have goods, or craftable or findable in loot, or bases that are flight reachable from manhattan"  example("true")
func GetCMs(webapp *web.Web, api *Api) *registry.Endpoint {
	return &registry.Endpoint{
		Url:     "GET " + ApiRoute + "/counter_measures",
		Handler: GetItemsT(webapp, api.app_data, api.app_data.Configs.CMs, api.app_data.Configs.FilterToUsefulCounterMeasures),
	}
}

// ShowAccount godoc
// @Summary      Getting list of CounterMeasure Market Goods
// @Tags         counter_measures
// @Accept       json
// @Produce      json
// @Param request body []string true "Array of counter nicknames as input, for example [ge_s_cm_01]"
// @Success      200  {array}  	MarketGoodResp
// @Router       /api/counter_measures/market_goods [post]
func PostCMsMarketGoods(webapp *web.Web, api *Api) *registry.Endpoint {
	return &registry.Endpoint{
		Url:     "POST " + ApiRoute + "/counter_measures/market_goods",
		Handler: PostItemsMarketGoodsT(webapp, api.app_data.Configs.CMs),
	}
}

// ShowAccount godoc
// @Summary      Getting list of CounterMeasure Tech compats
// @Tags         counter_measures
// @Accept       json
// @Produce      json
// @Param request body []string true "Array of counter measure nicknames as input"
// @Success      200  {array}  	TechCompatResp
// @Router       /api/counter_measures/tech_compats [post]
func PostCMsTechcompatibilities(webapp *web.Web, api *Api) *registry.Endpoint {
	return &registry.Endpoint{
		Url:     "POST " + ApiRoute + "/counter_measures/tech_compats",
		Handler: PostItemsTechCompatT(webapp, api.app_data.Configs.CMs),
	}
}
