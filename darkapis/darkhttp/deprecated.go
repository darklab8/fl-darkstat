package darkhttp

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/darklab8/fl-darkstat/configs/cfg"
	"github.com/darklab8/fl-darkstat/darkapis/darkhttp/apiutils"
	"github.com/darklab8/fl-darkstat/darkcore/web"
	"github.com/darklab8/fl-darkstat/darkcore/web/registry"
	"github.com/darklab8/fl-darkstat/darkstat/configs_export"
	"github.com/darklab8/fl-darkstat/darkstat/settings/logus"
	"github.com/darklab8/go-utils/utils/ptr"
)

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

func PostAmmoMarketGoods(webapp *web.Web, api *Api) *registry.Endpoint {
	return &registry.Endpoint{
		Url:     "POST " + ApiRoute + "/ammos/market_goods",
		Handler: PostItemsMarketGoodsT(webapp, api.app_data.Configs.Ammos),
	}
}

func PostAmmoTechcompatibilities(webapp *web.Web, api *Api) *registry.Endpoint {
	return &registry.Endpoint{
		Url:     "POST " + ApiRoute + "/ammos/tech_compats",
		Handler: PostItemsTechCompatT(webapp, api.app_data.Configs.Ammos),
	}
}

func PostCMsMarketGoods(webapp *web.Web, api *Api) *registry.Endpoint {
	return &registry.Endpoint{
		Url:     "POST " + ApiRoute + "/counter_measures/market_goods",
		Handler: PostItemsMarketGoodsT(webapp, api.app_data.Configs.CMs),
	}
}

func PostCMsTechcompatibilities(webapp *web.Web, api *Api) *registry.Endpoint {
	return &registry.Endpoint{
		Url:     "POST " + ApiRoute + "/counter_measures/tech_compats",
		Handler: PostItemsTechCompatT(webapp, api.app_data.Configs.CMs),
	}
}

func PostCommodityMarketGoods(webapp *web.Web, api *Api) *registry.Endpoint {
	return &registry.Endpoint{
		Url:     "POST " + ApiRoute + "/commodities/market_goods",
		Handler: PostItemsMarketGoodsT(webapp, api.app_data.Configs.Commodities),
	}
}

func PostEnginesMarketGoods(webapp *web.Web, api *Api) *registry.Endpoint {
	return &registry.Endpoint{
		Url:     "POST " + ApiRoute + "/engines/market_goods",
		Handler: PostItemsMarketGoodsT(webapp, api.app_data.Configs.Engines),
	}
}

func PostEnginesTechcompatibilities(webapp *web.Web, api *Api) *registry.Endpoint {
	return &registry.Endpoint{
		Url:     "POST " + ApiRoute + "/engines/tech_compats",
		Handler: PostItemsTechCompatT(webapp, api.app_data.Configs.Engines),
	}
}

func PostGunsMarketGoods(webapp *web.Web, api *Api) *registry.Endpoint {
	return &registry.Endpoint{
		Url:     "POST " + ApiRoute + "/guns/market_goods",
		Handler: PostItemsMarketGoodsT(webapp, api.app_data.Configs.Guns),
	}
}

func PostGunsTechcompatibilities(webapp *web.Web, api *Api) *registry.Endpoint {
	return &registry.Endpoint{
		Url:     "POST " + ApiRoute + "/guns/tech_compats",
		Handler: PostItemsTechCompatT(webapp, api.app_data.Configs.Guns),
	}
}

func PostMissilesMarketGoods(webapp *web.Web, api *Api) *registry.Endpoint {
	return &registry.Endpoint{
		Url:     "POST " + ApiRoute + "/missiles/market_goods",
		Handler: PostItemsMarketGoodsT(webapp, api.app_data.Configs.Missiles),
	}
}

func PostMissilesTechcompatibilities(webapp *web.Web, api *Api) *registry.Endpoint {
	return &registry.Endpoint{
		Url:     "POST " + ApiRoute + "/missiles/tech_compats",
		Handler: PostItemsTechCompatT(webapp, api.app_data.Configs.Missiles),
	}
}

func PostMinesMarketGoods(webapp *web.Web, api *Api) *registry.Endpoint {
	return &registry.Endpoint{
		Url:     "POST " + ApiRoute + "/mines/market_goods",
		Handler: PostItemsMarketGoodsT(webapp, api.app_data.Configs.Mines),
	}
}

func PostMinesTechcompatibilities(webapp *web.Web, api *Api) *registry.Endpoint {
	return &registry.Endpoint{
		Url:     "POST " + ApiRoute + "/mines/tech_compats",
		Handler: PostItemsTechCompatT(webapp, api.app_data.Configs.Mines),
	}
}

func PostScannersMarketGoods(webapp *web.Web, api *Api) *registry.Endpoint {
	return &registry.Endpoint{
		Url:     "POST " + ApiRoute + "/scanners/market_goods",
		Handler: PostItemsMarketGoodsT(webapp, api.app_data.Configs.Scanners),
	}
}

func PostScannersTechcompatibilities(webapp *web.Web, api *Api) *registry.Endpoint {
	return &registry.Endpoint{
		Url:     "POST " + ApiRoute + "/scanners/tech_compats",
		Handler: PostItemsTechCompatT(webapp, api.app_data.Configs.Scanners),
	}
}

func PostShieldsMarketGoods(webapp *web.Web, api *Api) *registry.Endpoint {
	return &registry.Endpoint{
		Url:     "POST " + ApiRoute + "/shields/market_goods",
		Handler: PostItemsMarketGoodsT(webapp, api.app_data.Configs.Shields),
	}
}

func PostShieldsTechcompatibilities(webapp *web.Web, api *Api) *registry.Endpoint {
	return &registry.Endpoint{
		Url:     "POST " + ApiRoute + "/shields/tech_compats",
		Handler: PostItemsTechCompatT(webapp, api.app_data.Configs.Shields),
	}
}

func PostShipMarketGoods(webapp *web.Web, api *Api) *registry.Endpoint {
	return &registry.Endpoint{
		Url:     "POST " + ApiRoute + "/ships/market_goods",
		Handler: PostItemsMarketGoodsT(webapp, api.app_data.Configs.Ships),
	}
}

func PostShipTechcompatibilities(webapp *web.Web, api *Api) *registry.Endpoint {
	return &registry.Endpoint{
		Url:     "POST " + ApiRoute + "/ships/tech_compats",
		Handler: PostItemsTechCompatT(webapp, api.app_data.Configs.Ships),
	}
}

func PostThrustersMarketGoods(webapp *web.Web, api *Api) *registry.Endpoint {
	return &registry.Endpoint{
		Url:     "POST " + ApiRoute + "/thrusters/market_goods",
		Handler: PostItemsMarketGoodsT(webapp, api.app_data.Configs.Thrusters),
	}
}

func PostThrustersTechcompatibilities(webapp *web.Web, api *Api) *registry.Endpoint {
	return &registry.Endpoint{
		Url:     "POST " + ApiRoute + "/thrusters/tech_compats",
		Handler: PostItemsTechCompatT(webapp, api.app_data.Configs.Thrusters),
	}
}

func PostTractorMarketGoods(webapp *web.Web, api *Api) *registry.Endpoint {
	return &registry.Endpoint{
		Url:     "POST " + ApiRoute + "/tractors/market_goods",
		Handler: PostItemsMarketGoodsT(webapp, api.app_data.Configs.Tractors),
	}
}
