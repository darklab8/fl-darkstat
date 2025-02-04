package darkapi

import (
	"github.com/darklab8/fl-darkstat/darkcore/web"
	"github.com/darklab8/fl-darkstat/darkcore/web/registry"
)

// ShowAccount godoc
// @Summary      Getting list of Shields
// @Tags         shields
// @Accept       json
// @Produce      json
// @Success      200  {array}  	configs_export.Shield
// @Router       /api/shields [get]
// @Param        filter_to_useful    query     string  false  "filter items only to useful, usually they are sold, or have goods, or craftable or findable in loot, or bases that are flight reachable from manhattan"
func GetShields(webapp *web.Web, api *Api) *registry.Endpoint {
	return &registry.Endpoint{
		Url:     "GET " + ApiRoute + "/shields",
		Handler: GetItemsT(webapp, api.app_data.Configs.Shields, api.app_data.Configs.FilterToUsefulShields),
	}
}

// ShowAccount godoc
// @Summary      Getting list of Shields Market Goods
// @Tags         shields
// @Accept       json
// @Produce      json
// @Param request body []string true "Array of ammo nicknames as input, for example [ai_shield_hf]"
// @Success      200  {array}  	MarketGoodResp
// @Router       /api/shields/market_goods [post]
func PostShieldsMarketGoods(webapp *web.Web, api *Api) *registry.Endpoint {
	return &registry.Endpoint{
		Url:     "POST " + ApiRoute + "/shields/market_goods",
		Handler: PostItemsMarketGoodsT(webapp, api.app_data.Configs.Shields),
	}
}

// ShowAccount godoc
// @Summary      Getting list of Shields Tech compats
// @Tags         shields
// @Accept       json
// @Produce      json
// @Param request body []string true "Array of ammo nicknames as input"
// @Success      200  {array}  	TechCompatResp
// @Router       /api/shields/tech_compats [post]
func PostShieldsTechcompatibilities(webapp *web.Web, api *Api) *registry.Endpoint {
	return &registry.Endpoint{
		Url:     "POST " + ApiRoute + "/shields/tech_compats",
		Handler: PostItemsTechCompatT(webapp, api.app_data.Configs.Shields),
	}
}
