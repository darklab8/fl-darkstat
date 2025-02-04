package darkgrpc

import (
	"context"
	"fmt"
	"log"
	"net"

	pb "github.com/darklab8/fl-darkstat/darkproto"
	"github.com/darklab8/fl-darkstat/darkstat/appdata"
	"google.golang.org/grpc"
)

// server is used to implement helloworld.GreeterServer.
type Server struct {
	pb.UnimplementedDarkGRpcServer
	app_data *appdata.AppData
}

func NewServer(app_data *appdata.AppData) *Server {
	return &Server{app_data: app_data}
}

// SayHello implements helloworld.GreeterServer
func (s *Server) GetBases(_ context.Context, in *pb.Empty) (*pb.GetBasesReply, error) {
	var bases []*pb.Base
	for _, base := range s.app_data.Configs.Bases {
		item := &pb.Base{
			Name:               base.Name,
			Archetypes:         base.Archetypes,
			Nickname:           string(base.Nickname),
			NicknameHash:       int64(base.NicknameHash),
			FactionName:        base.FactionName,
			System:             base.System,
			SystemNickname:     base.SystemNickname,
			SystemNicknameHash: int64(base.SystemNicknameHash),
			Region:             base.Region,
			StridName:          int32(base.StridName),
			InfocardID:         int32(base.InfocardID),
			File:               base.File.ToString(),
			BGCSBaseRunBy:      base.BGCS_base_run_by,
			Pos: &pb.Pos{
				X: base.Pos.X,
				Y: base.Pos.Y,
				Z: base.Pos.Z,
			},
			SectorCoord:            base.SectorCoord,
			IsTransportUnreachable: base.IsTransportUnreachable,
			Reachable:              base.Reachable,
			IsPob:                  base.IsPob,
		}

		bases = append(bases, item)
	}
	return &pb.GetBasesReply{Bases: bases}, nil
}

const Port = 50051

func (r *Server) Serve() {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", Port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	pb.RegisterDarkGRpcServer(s, r)
	log.Printf("server listening at %v", lis.Addr())
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
	fmt.Println("grpc server is launched")
}
