package darkhttp

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"

	"github.com/darklab8/fl-darkstat/configs/cfg"
	"github.com/darklab8/fl-darkstat/darkapis/darkhttp/apiutils"
	"github.com/darklab8/fl-darkstat/darkcore/web"
	"github.com/darklab8/fl-darkstat/darkstat/configs_export"
	"github.com/darklab8/fl-darkstat/darkstat/settings/logus"
	"github.com/darklab8/go-utils/utils/ptr"
)

type MarketGoodResp struct {
	MarketGoods []*configs_export.MarketGood `json:"market_goods"`
	Nickname    string                       `json:"nickname"  validate:"required"`
	Error       *string                      `json:"error,omitempty"`
}

type TechCompatResp struct {
	TechCompat *configs_export.DiscoveryTechCompat `json:"tech_compat"`
	Nickname   string                              `json:"nickname"  validate:"required"`
	Error      *string                             `json:"error,omitempty"`
}

func GetItemsT[T Nicknamable](webapp *web.Web, items []T, filter func(items []T) []T) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		if webapp.AppDataMutex != nil {
			webapp.AppDataMutex.Lock()
			defer webapp.AppDataMutex.Unlock()
		}

		param1 := r.URL.Query().Get("filter_to_useful")
		var result []T
		if param1 == "true" {
			result = filter(items)
		} else {
			result = items
		}

		apiutils.ReturnJson(&w, result)
	}
}

type Nicknamable interface {
	GetNickname() string
}

type Marketable interface {
	Nicknamable
	GetBases() map[cfg.BaseUniNick]*configs_export.MarketGood
}

func PostItemsMarketGoodsT[T Marketable](webapp *web.Web, items []T) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		if webapp.AppDataMutex != nil {
			webapp.AppDataMutex.Lock()
			defer webapp.AppDataMutex.Unlock()
		}

		var nicknames []string
		body, err := io.ReadAll(r.Body)
		if logus.Log.CheckError(err, "failed to read body") {
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprintf(w, "err to ready body")
			return
		}
		json.Unmarshal(body, &nicknames)
		if len(nicknames) == 0 {
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprintf(w, "input at least some nicknames into request body")
			return
		}

		var market_good_answers []*MarketGoodResp

		items_by_nick := make(map[string]T)
		for _, item := range items {
			items_by_nick[string(item.GetNickname())] = item
		}

		for _, input_nickname := range nicknames {
			answer := &MarketGoodResp{Nickname: string(input_nickname)}
			if item, ok := items_by_nick[input_nickname]; ok {
				for _, good := range item.GetBases() {
					answer.MarketGoods = append(answer.MarketGoods, good)
				}
			} else {
				answer.Error = ptr.Ptr("not existing nickname")
			}
			market_good_answers = append(market_good_answers, answer)

		}
		apiutils.ReturnJson(&w, market_good_answers)
	}
}

type TechCompatable interface {
	Nicknamable
	GetDiscoveryTechCompat() *configs_export.DiscoveryTechCompat
}

func PostItemsTechCompatT[T TechCompatable](webapp *web.Web, items []T) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		if webapp.AppDataMutex != nil {
			webapp.AppDataMutex.Lock()
			defer webapp.AppDataMutex.Unlock()
		}

		var nicknames []string
		body, err := io.ReadAll(r.Body)
		if logus.Log.CheckError(err, "failed to read body") {
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprintf(w, "err to ready body")
			return
		}
		json.Unmarshal(body, &nicknames)
		if len(nicknames) == 0 {
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprintf(w, "input at least some base nicknames into request body")
			return
		}

		var market_good_answers []*TechCompatResp

		ships_by_nick := make(map[string]T)
		for _, item := range items {
			ships_by_nick[string(item.GetNickname())] = item
		}

		for _, input_nickname := range nicknames {
			answer := &TechCompatResp{Nickname: string(input_nickname)}
			if ship, ok := ships_by_nick[input_nickname]; ok {
				answer.TechCompat = ship.GetDiscoveryTechCompat()
			} else {
				answer.Error = ptr.Ptr("not existing nickname")
			}
			market_good_answers = append(market_good_answers, answer)

		}
		apiutils.ReturnJson(&w, market_good_answers)
	}
}

func ReadJsonInput[T any](w http.ResponseWriter, r *http.Request, data *T) error {
	// Try to decode the request body into the struct. If there is an error,
	// respond to the client with the error message and a 400 status code.
	err := json.NewDecoder(r.Body).Decode(data)
	if err != nil {
		return errors.New("failed to read body")
	}
	return nil
}
