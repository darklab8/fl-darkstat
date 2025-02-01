package darkapi

import (
	"github.com/darklab8/fl-darkstat/darkcore/web"
	"github.com/darklab8/fl-darkstat/darkcore/web/registry"
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
		Url:     "GET " + ApiRoute + "/ammos",
		Handler: GetItemsT(webapp, api.app_data.Configs.Ammos),
	}
}

// ShowAccount godoc
// @Summary      Getting list of Ammo Market Goods
// @Tags         ammos
// @Accept       json
// @Produce      json
// @Param request body []string true "Array of ammo nicknames as input, for example [dsy_annihilator_torpedo_ammo]"
// @Success      200  {array}  	MarketGoodResp
// @Router       /api/ammos/market_goods [post]
func PostAmmoMarketGoods(webapp *web.Web, api *Api) *registry.Endpoint {
	return &registry.Endpoint{
		Url:     "POST " + ApiRoute + "/ammos/market_goods",
		Handler: PostItemsMarketGoodsT(webapp, api.app_data.Configs.Ammos),
	}
}

// ShowAccount godoc
// @Summary      Getting list of Ammos Tech compats
// @Tags         ammos
// @Accept       json
// @Produce      json
// @Param request body []string true "Array of ammo nicknames as input"
// @Success      200  {array}  	TechCompatResp
// @Router       /api/ammos/tech_compats [post]
func PostAmmoTechcompatibilities(webapp *web.Web, api *Api) *registry.Endpoint {
	return &registry.Endpoint{
		Url:     "POST " + ApiRoute + "/ammos/tech_compats",
		Handler: PostItemsTechCompatT(webapp, api.app_data.Configs.Ammos),
	}
}
