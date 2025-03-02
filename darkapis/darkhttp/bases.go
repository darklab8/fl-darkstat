package darkhttp

import (
	"net/http"

	"github.com/darklab8/fl-darkstat/darkapis/darkgrpc"
	pb "github.com/darklab8/fl-darkstat/darkapis/darkgrpc/statproto"
	"github.com/darklab8/fl-darkstat/darkapis/darkhttp/apiutils"
	"github.com/darklab8/fl-darkstat/darkcore/web"
	"github.com/darklab8/fl-darkstat/darkcore/web/registry"
	"github.com/darklab8/fl-darkstat/darkstat/configs_export"
)

type Base struct {
	*configs_export.Base
	MarketGoods []*configs_export.MarketGood `json:"market_goods"`
}

func GetBasesInput(w http.ResponseWriter, r *http.Request) (*pb.GetBasesInput, error) {
	var in *pb.GetBasesInput = &pb.GetBasesInput{}
	if err := ReadJsonInput(w, r, &in); err != nil && r.Method == "POST" {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return in, err
	}

	if r.URL.Query().Get("filter_to_useful") == "true" {
		in.FilterToUseful = true
	}
	if r.URL.Query().Get("include_market_goods") == "true" {
		in.IncludeMarketGoods = true
	}
	return in, nil
}

// ShowAccount godoc
// @Summary      Getting list of NPC Bases
// @Tags         bases
// @Accept       json
// @Produce      json
// @Success      200  {array}  	darkhttp.Base
// @Router       /api/npc_bases [post]
// @Param request body pb.GetBasesInput true "input variables"
// @Description  include_market_goods: "insert 'true' if wish to include market goods under 'market goods' key or not. Such data can add a lot of extra weight"
// @Description  filter_to_useful: Apply filtering same as darkstat does by default for its tab. Usually means showing only items that can be bought/crafted/or found
// @Description  filter_nicknames: filters by item nicknames (in those case by base nicknames)
// @Description  filter_market_good_category: filters market goods to specific category. valid categories are written in market goods in same named attribute. for example 'commodity'
func GetBases(webapp *web.Web, api *Api) *registry.Endpoint {
	return &registry.Endpoint{
		Url: "" + ApiRoute + "/npc_bases",
		Handler: func(w http.ResponseWriter, r *http.Request) {
			if webapp.AppDataMutex != nil {
				webapp.AppDataMutex.Lock()
				defer webapp.AppDataMutex.Unlock()
			}

			var in *pb.GetBasesInput
			in, err := GetBasesInput(w, r)
			if err != nil {
				return
			}

			var output []*pb.Base = darkgrpc.GetBasesNpc(api.app_data, in)
			apiutils.ReturnJson(&w, output)
		},
	}

}

// ShowAccount godoc
// @Summary      Getting list of Mining Operations
// @Tags         bases
// @Accept       json
// @Produce      json
// @Success      200  {array}  	darkhttp.Base
// @Router       /api/mining_operations [post]
// @Param request body pb.GetBasesInput true "input variables"
// @Description  include_market_goods: "insert 'true' if wish to include market goods under 'market goods' key or not. Such data can add a lot of extra weight"
// @Description  filter_to_useful: Apply filtering same as darkstat does by default for its tab. Usually means showing only items that can be bought/crafted/or found
// @Description  filter_nicknames: filters by item nicknames (in those case by base nicknames)
// @Description  filter_market_good_category: filters market goods to specific category. valid categories are written in market goods in same named attribute. for example 'commodity'
func GetOreFields(webapp *web.Web, api *Api) *registry.Endpoint {
	return &registry.Endpoint{
		Url: "" + ApiRoute + "/mining_operations",
		Handler: func(w http.ResponseWriter, r *http.Request) {
			if webapp.AppDataMutex != nil {
				webapp.AppDataMutex.Lock()
				defer webapp.AppDataMutex.Unlock()
			}
			var in *pb.GetBasesInput
			in, err := GetBasesInput(w, r)
			if err != nil {
				return
			}

			var result []*configs_export.Base
			if in.FilterToUseful {
				result = configs_export.FitlerToUsefulOres(api.app_data.Configs.MiningOperations)
			} else {
				result = api.app_data.Configs.MiningOperations
			}
			result = darkgrpc.FilterNicknames(in.FilterNicknames, result)

			var output []*Base
			for _, item := range result {
				answer := &Base{
					Base: item,
				}
				if in.IncludeMarketGoods {
					for _, good := range darkgrpc.FilterMarketGoodCategory(in.FilterMarketGoodCategory, item.MarketGoodsPerNick) {
						answer.MarketGoods = append(answer.MarketGoods, good)
					}
				}
				output = append(output, answer)
			}
			apiutils.ReturnJson(&w, output)
		},
	}
}

// ShowAccount godoc
// @Summary      Getting list of Player Owned Bases in Bases format. Lists only pobs that have known position coordinates
// @Tags         bases
// @Accept       json
// @Produce      json
// @Success      200  {array}  	darkhttp.Base
// @Router       /api/pobs/bases [post]
// @Param request body pb.GetBasesInput true "input variables"
// @Description  include_market_goods: "insert 'true' if wish to include market goods under 'market goods' key or not. Such data can add a lot of extra weight"
// @Description  filter_to_useful: Apply filtering same as darkstat does by default for its tab. Usually means showing only items that can be bought/crafted/or found
// @Description  filter_nicknames: filters by item nicknames (in those case by base nicknames)
// @Description  filter_market_good_category: filters market goods to specific category. valid categories are written in market goods in same named attribute. for example 'commodity'
func GetPoBBases(webapp *web.Web, api *Api) *registry.Endpoint {
	return &registry.Endpoint{
		Url: "" + ApiRoute + "/pobs/bases",
		Handler: func(w http.ResponseWriter, r *http.Request) {
			if webapp.AppDataMutex != nil {
				webapp.AppDataMutex.Lock()
				defer webapp.AppDataMutex.Unlock()
			}
			var in *pb.GetBasesInput
			in, err := GetBasesInput(w, r)
			if err != nil {
				return
			}

			var result []*configs_export.Base = api.app_data.Configs.PoBsToBases(api.app_data.Configs.PoBs)
			result = darkgrpc.FilterNicknames(in.FilterNicknames, result)

			var output []*Base
			for _, item := range result {
				answer := &Base{
					Base: item,
				}
				if in.IncludeMarketGoods {
					for _, good := range darkgrpc.FilterMarketGoodCategory(in.FilterMarketGoodCategory, item.MarketGoodsPerNick) {
						answer.MarketGoods = append(answer.MarketGoods, good)
					}
				}
				output = append(output, answer)
			}

			apiutils.ReturnJson(&w, output)
		},
	}
}
