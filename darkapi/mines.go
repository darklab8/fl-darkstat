package darkapi

import (
	"github.com/darklab8/fl-darkstat/darkcore/web"
	"github.com/darklab8/fl-darkstat/darkcore/web/registry"
)

// ShowAccount godoc
// @Summary      Getting list of Mines
// @Tags         mines
// @Accept       json
// @Produce      json
// @Success      200  {array}  	configs_export.Mine
// @Router       /api/mines [get]
// @Param        filter_to_useful    query     string  false  "filter items only to useful, usually they are sold, or have goods, or craftable or findable in loot, or bases that are flight reachable from manhattan"  example("true")
func GetMines(webapp *web.Web, api *Api) *registry.Endpoint {
	return &registry.Endpoint{
		Url:     "GET " + ApiRoute + "/mines",
		Handler: GetItemsT(webapp, api.app_data.Configs.Mines, api.app_data.Configs.FilterToUsefulMines),
	}
}

// ShowAccount godoc
// @Summary      Getting list of Mines Market Goods
// @Tags         mines
// @Accept       json
// @Produce      json
// @Param request body []string true "Array of mine nicknames as input, for example [mine02_mark02]"
// @Success      200  {array}  	MarketGoodResp
// @Router       /api/mines/market_goods [post]
func PostMinesMarketGoods(webapp *web.Web, api *Api) *registry.Endpoint {
	return &registry.Endpoint{
		Url:     "POST " + ApiRoute + "/mines/market_goods",
		Handler: PostItemsMarketGoodsT(webapp, api.app_data.Configs.Mines),
	}
}

// ShowAccount godoc
// @Summary      Getting list of Mines Tech compats
// @Tags         mines
// @Accept       json
// @Produce      json
// @Param request body []string true "Array of gun nicknames as input, for example [ai_bomber]" example("ai_bomber")
// @Success      200  {array}  	TechCompatResp
// @Router       /api/mines/tech_compats [post]
func PostMinesTechcompatibilities(webapp *web.Web, api *Api) *registry.Endpoint {
	return &registry.Endpoint{
		Url:     "POST " + ApiRoute + "/mines/tech_compats",
		Handler: PostItemsTechCompatT(webapp, api.app_data.Configs.Mines),
	}
}
