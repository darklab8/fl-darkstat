package darkgrpc

import (
	"context"

	"github.com/darklab8/fl-darkstat/darkapi"
	pb "github.com/darklab8/fl-darkstat/darkgrpc/statproto"
)

func (s *Server) GetHashes(_ context.Context, in *pb.Empty) (*pb.GetHashesReply, error) {
	if s.app_data != nil {
		s.app_data.Lock()
		defer s.app_data.Unlock()
	}

	answer := &pb.GetHashesReply{HashesByNick: make(map[string]*pb.Hash)}

	hashes := darkapi.GetHashesData(s.app_data)

	for key, hash := range hashes {
		answer.HashesByNick[key] = &pb.Hash{
			Int32:  hash.Int32,
			Uint32: hash.Uint32,
			Hex:    hash.Hex,
		}
	}

	return answer, nil
}
