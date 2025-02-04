package darkapi

import (
	"github.com/darklab8/fl-darkstat/darkcore/web"
	"github.com/darklab8/fl-darkstat/darkcore/web/registry"
)

// ShowAccount godoc
// @Summary      Getting list of Ships
// @Tags         ships
// @Accept       json
// @Produce      json
// @Success      200  {array}  	configs_export.Ship
// @Router       /api/ships [get]
// @Param        filter_to_useful    query     string  false  "filter items only to useful, usually they are sold, or have goods, or craftable or findable in loot, or bases that are flight reachable from manhattan"
func GetShips(webapp *web.Web, api *Api) *registry.Endpoint {
	return &registry.Endpoint{
		Url:     "GET " + ApiRoute + "/ships",
		Handler: GetItemsT(webapp, api.app_data.Configs.Ships, api.app_data.Configs.FilterToUsefulShips),
	}
}

// ShowAccount godoc
// @Summary      Getting list of Ship Market Goods
// @Tags         ships
// @Accept       json
// @Produce      json
// @Param request body []string true "Array of ship nicknames as input, for example [ai_bomber]"
// @Success      200  {array}  	MarketGoodResp
// @Router       /api/ships/market_goods [post]
func PostShipMarketGoods(webapp *web.Web, api *Api) *registry.Endpoint {
	return &registry.Endpoint{
		Url:     "POST " + ApiRoute + "/ships/market_goods",
		Handler: PostItemsMarketGoodsT(webapp, api.app_data.Configs.Ships),
	}
}

// ShowAccount godoc
// @Summary      Getting list of Ship Tech compats
// @Tags         ships
// @Accept       json
// @Produce      json
// @Param request body []string true "Array of ship nicknames as input, for example [ai_bomber]"
// @Success      200  {array}  	TechCompatResp
// @Router       /api/ships/tech_compats [post]
func PostShipTechcompatibilities(webapp *web.Web, api *Api) *registry.Endpoint {
	return &registry.Endpoint{
		Url:     "POST " + ApiRoute + "/ships/tech_compats",
		Handler: PostItemsTechCompatT(webapp, api.app_data.Configs.Ships),
	}
}
