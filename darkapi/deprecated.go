package darkapi

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/darklab8/fl-darkstat/configs/cfg"
	"github.com/darklab8/fl-darkstat/darkapi/apiutils"
	"github.com/darklab8/fl-darkstat/darkcore/web"
	"github.com/darklab8/fl-darkstat/darkcore/web/registry"
	"github.com/darklab8/fl-darkstat/darkstat/configs_export"
	"github.com/darklab8/fl-darkstat/darkstat/settings/logus"
	"github.com/darklab8/go-utils/utils/ptr"
)

// ShowAccount godoc
// @Summary      Getting list of NPC Bases Market Goods
// @Tags         deprecated
// @Accept       json
// @Produce      json
// @Param request body []cfg.BaseUniNick true "Array of npc base nicknames as input, for example [li01_01_base]"
// @Success      200  {array}  	MarketGoodResp
// @Router       /api/npc_bases/market_goods [post]
func PostBaseMarketGoods(webapp *web.Web, api *Api) *registry.Endpoint {
	return &registry.Endpoint{
		Url: "POST " + ApiRoute + "/npc_bases/market_goods",
		Handler: func(w http.ResponseWriter, r *http.Request) {
			if webapp.AppDataMutex != nil {
				webapp.AppDataMutex.Lock()
				defer webapp.AppDataMutex.Unlock()
			}

			var base_nicknames []cfg.BaseUniNick
			body, err := io.ReadAll(r.Body)
			if logus.Log.CheckError(err, "failed to read body") {
				w.WriteHeader(http.StatusBadRequest)
				fmt.Fprintf(w, "err to ready body")
				return
			}
			json.Unmarshal(body, &base_nicknames)
			if len(base_nicknames) == 0 {
				w.WriteHeader(http.StatusBadRequest)
				fmt.Fprintf(w, "input at least some base nicknames into request body")
				return
			}

			var base_market_goods []*MarketGoodResp

			bases_by_nick := make(map[string]*configs_export.Base)
			for _, base := range api.app_data.Configs.Bases {
				bases_by_nick[string(base.Nickname)] = base
			}

			for _, base_nickname := range base_nicknames {
				answer := &MarketGoodResp{Nickname: string(base_nickname)}
				if base, ok := bases_by_nick[base_nickname.ToStr()]; ok {
					for _, good := range base.MarketGoodsPerNick {
						answer.MarketGoods = append(answer.MarketGoods, good)
					}
				} else {
					answer.Error = ptr.Ptr("not existing base")
				}
				base_market_goods = append(base_market_goods, answer)

			}
			apiutils.ReturnJson(&w, base_market_goods)
		},
	}
}

// ShowAccount godoc
// @Summary      Getting list of Ammo Market Goods
// @Tags         deprecated
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
// @Tags         deprecated
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

// ShowAccount godoc
// @Summary      Getting list of CounterMeasure Market Goods
// @Tags         deprecated
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
// @Tags         deprecated
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

// ShowAccount godoc
// @Summary      Getting list of Commodities Market Goods
// @Tags         deprecated
// @Accept       json
// @Produce      json
// @Param request body []string true "Array of commodity nicknames as input, for example [commodity_military_salvage]"
// @Success      200  {array}  	MarketGoodResp
// @Router       /api/commodities/market_goods [post]
func PostCommodityMarketGoods(webapp *web.Web, api *Api) *registry.Endpoint {
	return &registry.Endpoint{
		Url:     "POST " + ApiRoute + "/commodities/market_goods",
		Handler: PostItemsMarketGoodsT(webapp, api.app_data.Configs.Commodities),
	}
}

// ShowAccount godoc
// @Summary      Getting list of Engine Market Goods
// @Tags         deprecated
// @Accept       json
// @Produce      json
// @Param request body []string true "Array of engines nicknames as input, for example [ge_kfr_engine_01_add]"
// @Success      200  {array}  	MarketGoodResp
// @Router       /api/engines/market_goods [post]
func PostEnginesMarketGoods(webapp *web.Web, api *Api) *registry.Endpoint {
	return &registry.Endpoint{
		Url:     "POST " + ApiRoute + "/engines/market_goods",
		Handler: PostItemsMarketGoodsT(webapp, api.app_data.Configs.Engines),
	}
}

// ShowAccount godoc
// @Summary      Getting list of Engine Tech compats
// @Tags         deprecated
// @Accept       json
// @Produce      json
// @Param request body []string true "Array of engines nicknames as input"
// @Success      200  {array}  	TechCompatResp
// @Router       /api/engines/tech_compats [post]
func PostEnginesTechcompatibilities(webapp *web.Web, api *Api) *registry.Endpoint {
	return &registry.Endpoint{
		Url:     "POST " + ApiRoute + "/engines/tech_compats",
		Handler: PostItemsTechCompatT(webapp, api.app_data.Configs.Engines),
	}
}

// ShowAccount godoc
// @Summary      Getting list of Guns Market Goods
// @Tags         deprecated
// @Accept       json
// @Produce      json
// @Param request body []string true "Array of ship nicknames as input, for example [ai_bomber]"
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
// @Tags         deprecated
// @Accept       json
// @Produce      json
// @Param request body []string true "Array of gun nicknames as input, for example [ai_bomber]"
// @Success      200  {array}  	TechCompatResp
// @Router       /api/guns/tech_compats [post]
func PostGunsTechcompatibilities(webapp *web.Web, api *Api) *registry.Endpoint {
	return &registry.Endpoint{
		Url:     "POST " + ApiRoute + "/guns/tech_compats",
		Handler: PostItemsTechCompatT(webapp, api.app_data.Configs.Guns),
	}
}

// ShowAccount godoc
// @Summary      Getting list of Missiles Market Goods
// @Tags         deprecated
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
// @Tags         deprecated
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

// ShowAccount godoc
// @Summary      Getting list of Mines Market Goods
// @Tags         deprecated
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
// @Tags         deprecated
// @Accept       json
// @Produce      json
// @Param request body []string true "Array of gun nicknames as input, for example [ai_bomber]"
// @Success      200  {array}  	TechCompatResp
// @Router       /api/mines/tech_compats [post]
func PostMinesTechcompatibilities(webapp *web.Web, api *Api) *registry.Endpoint {
	return &registry.Endpoint{
		Url:     "POST " + ApiRoute + "/mines/tech_compats",
		Handler: PostItemsTechCompatT(webapp, api.app_data.Configs.Mines),
	}
}

// ShowAccount godoc
// @Summary      Getting list of Scanners Market Goods
// @Tags         deprecated
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
// @Tags         deprecated
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

// ShowAccount godoc
// @Summary      Getting list of Shields Market Goods
// @Tags         deprecated
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
// @Tags         deprecated
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

// ShowAccount godoc
// @Summary      Getting list of Ship Market Goods
// @Tags         deprecated
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
// @Tags         deprecated
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

// ShowAccount godoc
// @Summary      Getting list of Thrusters Market Goods
// @Tags         deprecated
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
// @Tags         deprecated
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

// ShowAccount godoc
// @Summary      Getting list of Tractor Market Goods
// @Tags         deprecated
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
