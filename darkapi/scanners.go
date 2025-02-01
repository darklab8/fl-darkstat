package darkapi

import (
	"github.com/darklab8/fl-darkstat/darkcore/web"
	"github.com/darklab8/fl-darkstat/darkcore/web/registry"
)

// ShowAccount godoc
// @Summary      Getting list of Scanners
// @Tags         scanners
// @Accept       json
// @Produce      json
// @Success      200  {array}  	configs_export.Scanner
// @Router       /api/scanners [get]
func GetScanners(webapp *web.Web, api *Api) *registry.Endpoint {
	return &registry.Endpoint{
		Url:     "GET " + ApiRoute + "/scanners",
		Handler: GetItemsT(webapp, api.app_data.Configs.Scanners),
	}
}

// ShowAccount godoc
// @Summary      Getting list of Scanners Market Goods
// @Tags         scanners
// @Accept       json
// @Produce      json
// @Param request body []string true "Array of ammo nicknames as input, for example [dsy_annihilator_torpedo_ammo]"
// @Success      200  {array}  	MarketGoodResp
// @Router       /api/scanners/market_goods [post]
func PostScannersMarketGoods(webapp *web.Web, api *Api) *registry.Endpoint {
	return &registry.Endpoint{
		Url:     "POST " + ApiRoute + "/scanners/market_goods",
		Handler: PostItemsMarketGoodsT(webapp, api.app_data.Configs.Scanners),
	}
}

// ShowAccount godoc
// @Summary      Getting list of Scanners Tech compats
// @Tags         scanners
// @Accept       json
// @Produce      json
// @Param request body []string true "Array of ammo nicknames as input"
// @Success      200  {array}  	TechCompatResp
// @Router       /api/scanners/tech_compats [post]
func PostScannersTechcompatibilities(webapp *web.Web, api *Api) *registry.Endpoint {
	return &registry.Endpoint{
		Url:     "POST " + ApiRoute + "/scanners/tech_compats",
		Handler: PostItemsTechCompatT(webapp, api.app_data.Configs.Scanners),
	}
}
