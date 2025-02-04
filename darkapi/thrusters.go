package darkapi

import (
	"github.com/darklab8/fl-darkstat/darkcore/web"
	"github.com/darklab8/fl-darkstat/darkcore/web/registry"
)

// ShowAccount godoc
// @Summary      Getting list of Thrusters
// @Tags         thrusters
// @Accept       json
// @Produce      json
// @Success      200  {array}  	configs_export.Thruster
// @Router       /api/thrusters [get]
// @Param        filter_to_useful    query     string  false  "filter items only to useful, usually they are sold, or have goods, or craftable or findable in loot, or bases that are flight reachable from manhattan"
func GetThrusters(webapp *web.Web, api *Api) *registry.Endpoint {
	return &registry.Endpoint{
		Url:     "GET " + ApiRoute + "/thrusters",
		Handler: GetItemsT(webapp, api.app_data.Configs.Thrusters, api.app_data.Configs.FilterToUsefulThrusters),
	}
}

// ShowAccount godoc
// @Summary      Getting list of Thrusters Market Goods
// @Tags         thrusters
// @Accept       json
// @Produce      json
// @Param request body []string true "Array of thrusters nicknames as input, for example [dsy_thruster_bd]"
// @Success      200  {array}  	MarketGoodResp
// @Router       /api/thrusters/market_goods [post]
func PostThrustersMarketGoods(webapp *web.Web, api *Api) *registry.Endpoint {
	return &registry.Endpoint{
		Url:     "POST " + ApiRoute + "/thrusters/market_goods",
		Handler: PostItemsMarketGoodsT(webapp, api.app_data.Configs.Thrusters),
	}
}

// ShowAccount godoc
// @Summary      Getting list of Thrusters Tech compats
// @Tags         thrusters
// @Accept       json
// @Produce      json
// @Param request body []string true "Array of thrusters nicknames as input"
// @Success      200  {array}  	TechCompatResp
// @Router       /api/thrusters/tech_compats [post]
func PostThrustersTechcompatibilities(webapp *web.Web, api *Api) *registry.Endpoint {
	return &registry.Endpoint{
		Url:     "POST " + ApiRoute + "/thrusters/tech_compats",
		Handler: PostItemsTechCompatT(webapp, api.app_data.Configs.Thrusters),
	}
}
