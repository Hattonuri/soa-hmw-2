package grpc_server

import (
	"context"

	"github.com/hattonuri/soa-hmw-2/internal/config"
	"github.com/hattonuri/soa-hmw-2/internal/domain/usecases"
	"github.com/hattonuri/soa-hmw-2/internal/generated/proto/base"
	"google.golang.org/protobuf/types/known/emptypb"
)

type GRPCServer struct {
	base.UnimplementedMafiaServiceServer
	server *usecases.Server
}

func (s *GRPCServer) Join(in *base.Player, srv base.MafiaService_JoinServer) error {
	return s.server.Join(in, srv)
}

func (s *GRPCServer) ResponseEvent(ctx context.Context, event *base.Event) (*emptypb.Empty, error) {
	s.server.ResponseEvent(event)
	return &emptypb.Empty{}, nil
}

func (s *GRPCServer) GetPlayersList(ctx context.Context, room *base.Room) (*base.Players, error) {
	return s.server.GetPlayersList(room)
}

func NewServer(c *config.Server) *GRPCServer {
	return &GRPCServer{
		server: usecases.NewServer(c),
	}
}
