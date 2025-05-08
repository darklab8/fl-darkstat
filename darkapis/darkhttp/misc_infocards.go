package darkhttp

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/darklab8/fl-darkstat/darkapis/darkhttp/apiutils"
	"github.com/darklab8/fl-darkstat/darkcore/web"
	"github.com/darklab8/fl-darkstat/darkcore/web/registry"
	"github.com/darklab8/fl-darkstat/darkstat/appdata"
	"github.com/darklab8/fl-darkstat/darkstat/configs_export/infocarder"
	"github.com/darklab8/fl-darkstat/darkstat/settings/logus"
	"github.com/darklab8/go-utils/utils/ptr"
)

// ShowAccount godoc
// @Summary      Getting infocards
// @Tags         misc
// @Accept       json
// @Produce      json
// @Param request body []string true "Array of nicknames as input, for example [fc_or_gun01_mark02]"
// @Success      200  {array}  	InfocardResp
// @Router       /api/infocards [post]
func GetInfocards(webapp *web.Web, app_data *appdata.AppData, api *Api) *registry.Endpoint {
	return &registry.Endpoint{
		Url: "POST " + ApiRoute + "/infocards",
		Handler: func(w http.ResponseWriter, r *http.Request) {
			if webapp.AppDataMutex != nil {
				webapp.AppDataMutex.Lock()
				defer webapp.AppDataMutex.Unlock()
			}

			var nicknames []string
			body, err := io.ReadAll(r.Body)
			if logus.Log.CheckError(err, "failed to read body") {
				w.WriteHeader(http.StatusBadRequest)
				_, err = fmt.Fprintf(w, "err to ready body")
				Log.CheckError(err, "fprintf print error in infocards 1")
				return
			}
			err = json.Unmarshal(body, &nicknames)
			Log.CheckWarn(err, "failed to unparmshal input in get infocards")
			if len(nicknames) == 0 {
				w.WriteHeader(http.StatusBadRequest)
				_, err = fmt.Fprintf(w, "input at least some base nicknames into request body")
				Log.CheckError(err, "fprintf print error in infocards 2")
				return
			}

			var outputs []InfocardResp
			for _, nickname := range nicknames {
				if info, ok := app_data.Configs.Infocards[infocarder.InfocardKey(nickname)]; ok {
					outputs = append(outputs, InfocardResp{Infocard: &info})
				} else {
					outputs = append(outputs, InfocardResp{Error: ptr.Ptr("infocard is not found")})
				}
			}

			apiutils.ReturnJson(&w, outputs)
		},
	}
}

type InfocardResp struct {
	Infocard *infocarder.Infocard `json:"infocard,omitempty"`
	Error    *string              `json:"error,omitempty"`
}
