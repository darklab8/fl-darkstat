package darkgrpc

import (
	"context"

	pb "github.com/darklab8/fl-darkstat/darkgrpc/statproto"
	"github.com/darklab8/fl-darkstat/darkstat/configs_export"
	"github.com/darklab8/go-utils/utils/ptr"
)

func (s *Server) GetInfocards(_ context.Context, in *pb.GetInfocardsInput) (*pb.GetInfocardsReply, error) {
	if s.app_data != nil {
		s.app_data.Lock()
		defer s.app_data.Unlock()
	}

	var outputs []*pb.GetInfocardAnswer
	for _, nickname := range in.Nicknames {
		if info, ok := s.app_data.Configs.Infocards[configs_export.InfocardKey(nickname)]; ok {
			outputs = append(outputs, &pb.GetInfocardAnswer{Infocard: NewInfocard(info)})
		} else {
			outputs = append(outputs, &pb.GetInfocardAnswer{Error: ptr.Ptr("infocard is not found")})
		}
	}

	return &pb.GetInfocardsReply{Answers: outputs}, nil
}

func NewInfocard(info configs_export.Infocard) *pb.Infocard {
	result := &pb.Infocard{}

	for _, line := range info {
		line_to_add := &pb.InfocardLine{}

		for _, phrase := range line.Phrases {
			phrase_to_add := &pb.InfocardPhrase{
				Phrase: phrase.Phrase,
				Link:   phrase.Link,
				Bold:   phrase.Bold,
			}
			line_to_add.Phrases = append(line_to_add.Phrases, phrase_to_add)
		}
		result.Lines = append(result.Lines, line_to_add)
	}
	return result
}
