package darkgrpc

import (
	"context"

	pb "github.com/darklab8/fl-darkstat/darkapis/darkgrpc/statproto"
	"github.com/darklab8/fl-darkstat/darkstat/appdata"
)

func (s *Server) GetGraphPaths(_ context.Context, in *pb.GetGraphPathsInput) (*pb.GetGraphPathsReply, error) {
	if s.app_data != nil {
		s.app_data.RLock()
		defer s.app_data.RUnlock()
	}

	var input_queries []appdata.GraphPathReq
	for _, query := range in.Queries {
		input_queries = append(input_queries, appdata.GraphPathReq{
			From: query.From,
			To:   query.To,
		})
	}
	answers := s.app_data.GetGraphPaths(input_queries)
	var grpc_answers []*pb.GetGraphPathsAnswer

	for _, answer := range answers {
		grpc_answers = append(grpc_answers, &pb.GetGraphPathsAnswer{
			Route: NewGraphQuery(&answer.Query),
			Time:  NewGraphTimeAnswer(answer.Time),
			Error: answer.Error,
		})
	}

	return &pb.GetGraphPathsReply{Answers: grpc_answers}, nil
}

func NewGraphQuery(Query *appdata.GraphPathReq) *pb.GraphPathQuery {
	if Query == nil {
		return nil
	}

	return &pb.GraphPathQuery{
		From: Query.From,
		To:   Query.To,
	}
}

func NewGraphTimeAnswer(Time *appdata.GraphPathTime) *pb.GraphPathTime {
	if Time == nil {
		return nil
	}

	return &pb.GraphPathTime{
		Transport: NewInt64(Time.Transport),
		Frigate:   NewInt64(Time.Frigate),
		Freighter: NewInt64(Time.Freighter),
	}
}
