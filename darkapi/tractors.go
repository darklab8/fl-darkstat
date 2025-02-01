package darkapi

import (
	"github.com/darklab8/fl-darkstat/darkcore/web"
	"github.com/darklab8/fl-darkstat/darkcore/web/registry"
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
		Url:     "GET " + ApiRoute + "/tractors",
		Handler: GetItemsT(webapp, api.app_data.Configs.Tractors),
	}
}

// ShowAccount godoc
// @Summary      Getting list of Tractor Market Goods
// @Tags         tractors
// @Accept       json
// @Produce      json
// @Param request body []string true "Array of ship nicknames as input, for example [dsy_license_srp_28]"
// @Success      200  {array}  	MarketGoodResp
// @Router       /api/tractors/market_goods [post]
func PostTractorMarketGoods(webapp *web.Web, api *Api) *registry.Endpoint {
	return &registry.Endpoint{
		Url:     "POST " + ApiRoute + "/tractors/market_goods",
		Handler: PostItemsMarketGoodsT(webapp, api.app_data.Configs.Tractors),
	}
}
