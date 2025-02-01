package darkapi

import (
	"github.com/darklab8/fl-darkstat/darkcore/web"
	"github.com/darklab8/fl-darkstat/darkcore/web/registry"
)

// ShowAccount godoc
// @Summary      Getting list of Guns
// @Tags         guns
// @Accept       json
// @Produce      json
// @Success      200  {array}  	configs_export.Gun
// @Router       /api/guns [get]
// @Param        filter_to_useful    query     string  false  "filter items only to useful, usually they are sold, or have goods, or craftable or findable in loot, or bases that are flight reachable from manhattan"  example("true")
func GetGuns(webapp *web.Web, api *Api) *registry.Endpoint {
	return &registry.Endpoint{
		Url:     "GET " + ApiRoute + "/guns",
		Handler: GetItemsT(webapp, api.app_data, api.app_data.Configs.Guns, api.app_data.Configs.FilterToUsefulGun),
	}
}

// ShowAccount godoc
// @Summary      Getting list of Guns Market Goods
// @Tags         guns
// @Accept       json
// @Produce      json
// @Param request body []string true "Array of ship nicknames as input, for example [ai_bomber]" example("ai_bomber")
// @Success      200  {array}  	MarketGoodResp
// @Router       /api/guns/market_goods [post]
func PostGunsMarketGoods(webapp *web.Web, api *Api) *registry.Endpoint {
	return &registry.Endpoint{
		Url:     "POST " + ApiRoute + "/guns/market_goods",
		Handler: PostItemsMarketGoodsT(webapp, api.app_data.Configs.Guns),
	}
}

// ShowAccount godoc
// @Summary      Getting list of Guns Tech compats
// @Tags         guns
// @Accept       json
// @Produce      json
// @Param request body []string true "Array of gun nicknames as input, for example [ai_bomber]" example("ai_bomber")
// @Success      200  {array}  	TechCompatResp
// @Router       /api/guns/tech_compats [post]
func PostGunsTechcompatibilities(webapp *web.Web, api *Api) *registry.Endpoint {
	return &registry.Endpoint{
		Url:     "POST " + ApiRoute + "/guns/tech_compats",
		Handler: PostItemsTechCompatT(webapp, api.app_data.Configs.Guns),
	}
}

// ShowAccount godoc
// @Summary      Getting list of Missiles
// @Tags         guns
// @Accept       json
// @Produce      json
// @Success      200  {array}  	configs_export.Gun
// @Router       /api/missiles [get]
// @Param        filter_to_useful    query     string  false  "filter items only to useful, usually they are sold, or have goods, or craftable or findable in loot, or bases that are flight reachable from manhattan"  example("true")
func GetMissiles(webapp *web.Web, api *Api) *registry.Endpoint {
	return &registry.Endpoint{
		Url:     "GET " + ApiRoute + "/missiles",
		Handler: GetItemsT(webapp, api.app_data, api.app_data.Configs.Missiles, api.app_data.Configs.FilterToUsefulGun),
	}
}

// ShowAccount godoc
// @Summary      Getting list of Missiles Market Goods
// @Tags         guns
// @Accept       json
// @Produce      json
// @Param request body []string true "Array of ship nicknames as input, for example [fc_or_gun01_mark02]"
// @Success      200  {array}  	MarketGoodResp
// @Router       /api/missiles/market_goods [post]
func PostMissilesMarketGoods(webapp *web.Web, api *Api) *registry.Endpoint {
	return &registry.Endpoint{
		Url:     "POST " + ApiRoute + "/missiles/market_goods",
		Handler: PostItemsMarketGoodsT(webapp, api.app_data.Configs.Missiles),
	}
}

// ShowAccount godoc
// @Summary      Getting list of Missiles Tech compats
// @Tags         guns
// @Accept       json
// @Produce      json
// @Param request body []string true "Array of missile nicknames as input, for example [fc_or_gun01_mark02]"
// @Success      200  {array}  	TechCompatResp
// @Router       /api/missiles/tech_compats [post]
func PostMissilesTechcompatibilities(webapp *web.Web, api *Api) *registry.Endpoint {
	return &registry.Endpoint{
		Url:     "POST " + ApiRoute + "/missiles/tech_compats",
		Handler: PostItemsTechCompatT(webapp, api.app_data.Configs.Missiles),
	}
}
