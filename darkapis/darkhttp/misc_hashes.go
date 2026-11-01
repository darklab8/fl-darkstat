package darkhttp

import (
	"net/http"

	"github.com/darklab8/fl-darkstat/darkapis/darkhttp/apiutils"
	pb "github.com/darklab8/fl-darkstat/darkapis/darkhttp/pb"
	"github.com/darklab8/fl-darkstat/darkcore/web"
	"github.com/darklab8/fl-darkstat/darkcore/web/registry"
)

// ShowAccount godoc
// @Summary      Hashes
// @Tags         misc
// @Accept       json
// @Produce      json
// @Success      200  {object}  	pb.Hashes
// @Router       /api/hashes [post]
func GetHashes(webapp *web.Web, api *Api) *registry.Endpoint {
	return &registry.Endpoint{
		Url: "" + ApiRoute + "/hashes",
		Handler: func(w http.ResponseWriter, r *http.Request) {
			if webapp.AppDataMutex != nil {
				webapp.AppDataMutex.RLock()
				defer webapp.AppDataMutex.RUnlock()
			}
			hashes := pb.GetHashesData(api.app_data)
			apiutils.ReturnJson(&w, pb.Hashes{HashesByNick: hashes})
		},
	}
}
func (c *HttpClient) GetHashes() (pb.Hashes, error) {
	return make_request[EmptyInput, pb.Hashes](c, ""+ApiRoute+"/hashes", EmptyInput{})
}
