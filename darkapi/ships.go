package darkapi

import (
	"github.com/darklab8/fl-darkstat/darkcore/web"
	"github.com/darklab8/fl-darkstat/darkcore/web/registry"
	"github.com/darklab8/fl-darkstat/darkstat/configs_export"
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
		Url:     "GET " + ApiRoute + "/ships",
		Handler: GetItemsT(webapp, api.app_data.Configs.Ships),
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
		Url:     "POST " + ApiRoute + "/ships/market_goods",
		Handler: PostItemsMarketGoodsT(webapp, api.app_data.Configs.Ships),
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
		Url:     "POST " + ApiRoute + "/ships/tech_compats",
		Handler: PostItemsTechCompatT(webapp, api.app_data.Configs.Ships),
	}
}
