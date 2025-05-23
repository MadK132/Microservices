package server

import (
	"context"
	"log"
	"net"

	statisticspb "github.com/recktt77/proto-definitions/gen/statistics"
	"github.com/recktt77/statistics_service/internal/repository"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

)

type StatsServer struct {
	statisticspb.UnimplementedStatisticsServiceServer
	repo repository.StatsRepository
}

func NewGRPCServer(repo repository.StatsRepository) *StatsServer {
	return &StatsServer{repo: repo}
}

func (s *StatsServer) GetStatistics(ctx context.Context, _ *statisticspb.Empty) (*statisticspb.GetStatisticsResponse, error) {
	data, err := s.repo.GetAll()
	if err != nil {
		return nil, err
	}

	resp := &statisticspb.GetStatisticsResponse{
		Stats: make(map[string]*statisticspb.ActionCount),
	}

	for source, actions := range data {
		resp.Stats[source] = &statisticspb.ActionCount{
			Actions: make(map[string]int32),
		}
		for action, count := range actions {
			resp.Stats[source].Actions[action] = int32(count)
		}
	}

	return resp, nil
}


func (s *StatsServer) Run(port string) error {
	lis, err := net.Listen("tcp", port)
	if err != nil {
		return err
	}

	grpcServer := grpc.NewServer()
	statisticspb.RegisterStatisticsServiceServer(grpcServer, s)

	reflection.Register(grpcServer)

	log.Println("gRPC StatisticsService listening on", port)
	return grpcServer.Serve(lis)

}
