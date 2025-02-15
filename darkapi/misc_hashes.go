package darkapi

import (
	"net/http"
	"runtime"
	"strings"
	"sync"

	"github.com/darklab8/fl-darkstat/configs/configs_mapped/freelancer_mapped/data_mapped/initialworld/flhash"
	"github.com/darklab8/fl-darkstat/configs/configs_mapped/parserutils/filefind"
	"github.com/darklab8/fl-darkstat/configs/configs_mapped/parserutils/iniload"
	"github.com/darklab8/fl-darkstat/darkapi/apiutils"
	"github.com/darklab8/fl-darkstat/darkcore/web"
	"github.com/darklab8/fl-darkstat/darkcore/web/registry"
	"github.com/darklab8/fl-darkstat/darkstat/appdata"
	"github.com/darklab8/fl-darkstat/darkstat/settings"
)

type Hash struct {
	Int32  int32  `json:"int32"  validate:"required"`
	Uint32 uint32 `json:"uint32"  validate:"required"`
	Hex    string `json:"hex"  validate:"required"`
}

type Hashes struct {
	HashesByNick map[string]Hash `json:"hashes_by_nick"  validate:"required"`
}

var hashes map[string]Hash

// ShowAccount godoc
// @Summary      Hashes
// @Tags         misc
// @Accept       json
// @Produce      json
// @Success      200  {object}  	Hashes
// @Router       /api/hashes [get]
func GetHashes(webapp *web.Web, api *Api) *registry.Endpoint {
	return &registry.Endpoint{
		Url: "GET " + ApiRoute + "/hashes",
		Handler: func(w http.ResponseWriter, r *http.Request) {
			if webapp.AppDataMutex != nil {
				webapp.AppDataMutex.Lock()
				defer webapp.AppDataMutex.Unlock()
			}

			hashes = GetHashesData(api.app_data)

			apiutils.ReturnJson(&w, Hashes{HashesByNick: hashes})
		},
	}
}

func GetHashesData(app_data *appdata.AppData) map[string]Hash {
	if hashes != nil {
		return hashes
	}

	hashes = make(map[string]Hash)

	filesystem := filefind.FindConfigs(settings.Env.FreelancerFolder)

	var wg sync.WaitGroup
	var mu sync.Mutex
	i := 0
	for filepath, file := range filesystem.Hashmap {
		if strings.Contains(filepath.Base().ToString(), "ini") {
			wg.Add(1)
			func(file *iniload.IniLoader) {
				file.Scan()
				for _, section := range file.Sections {
					if value, ok := section.ParamMap["nickname"]; ok {
						nickname := value[0].First.AsString()
						hash := flhash.HashNickname(nickname)
						mu.Lock()
						hashes[nickname] = Hash{
							Int32:  int32(hash),
							Uint32: uint32(hash),
							Hex:    hash.ToHexStr(),
						}
						mu.Unlock()
					}
				}
				wg.Done()
			}(iniload.NewLoader(file))
			i++
			if i%500 == 0 {
				runtime.GC()
			}
		}
	}
	wg.Wait()
	runtime.GC()

	for _, group := range app_data.Mapped.InitialWorld.Groups {
		var nickname string = group.Nickname.Get()
		hash := flhash.HashFaction(nickname)
		hashes[nickname] = Hash{
			Int32:  int32(hash),
			Uint32: uint32(hash),
			Hex:    hash.ToHexStr(),
		}
	}
	return hashes
}
