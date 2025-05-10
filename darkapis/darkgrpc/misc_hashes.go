package darkgrpc

import (
	"context"

	"runtime"
	"strings"
	"sync"

	"github.com/darklab8/fl-darkstat/configs/configs_mapped/freelancer_mapped/data_mapped/initialworld"
	"github.com/darklab8/fl-darkstat/configs/configs_mapped/freelancer_mapped/data_mapped/initialworld/flhash"
	"github.com/darklab8/fl-darkstat/configs/configs_mapped/parserutils/filefind"
	"github.com/darklab8/fl-darkstat/configs/configs_mapped/parserutils/iniload"
	"github.com/darklab8/fl-darkstat/configs/configs_settings"
	pb "github.com/darklab8/fl-darkstat/darkapis/darkgrpc/statproto"
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

	filesystem = filefind.FindConfigs(configs_settings.Env.FreelancerFolder)
	fileref := filesystem.GetFile(initialworld.FILENAME)
	InitialWorld := initialworld.Read(iniload.NewLoader(fileref).Scan())

	for _, group := range InitialWorld.Groups {
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

func (s *Server) GetHashes(_ context.Context, in *pb.Empty) (*pb.GetHashesReply, error) {
	if s.app_data != nil {
		s.app_data.RLock()
		defer s.app_data.RUnlock()
	}

	answer := &pb.GetHashesReply{HashesByNick: make(map[string]*pb.Hash)}

	hashes := GetHashesData(s.app_data)

	for key, hash := range hashes {
		answer.HashesByNick[key] = &pb.Hash{
			Int32:  hash.Int32,
			Uint32: hash.Uint32,
			Hex:    hash.Hex,
		}
	}

	return answer, nil
}
