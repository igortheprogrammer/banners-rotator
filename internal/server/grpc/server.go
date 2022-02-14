//go:generate mkdir -p ../bannersrotatorpb
//go:generate protoc -I ../../../api/ --go_out=../bannersrotatorpb/ ../../../api/BannersRotatorService.proto
//go:generate protoc -I ../../../api/ --go-grpc_out=../bannersrotatorpb/ ../../../api/BannersRotatorService.proto

package internalgrpc

import (
	"banners-rotator/internal/rotator"
	gw "banners-rotator/internal/server/bannersrotatorpb"
	"context"
	"errors"
	"fmt"
	"net"

	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	grpc_zap "github.com/grpc-ecosystem/go-grpc-middleware/logging/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var ErrBadRequest = errors.New("bad request")

type Server struct {
	app      rotator.App
	logger   rotator.Logger
	srv      *grpc.Server
	endpoint string
	gw.UnimplementedBannersRotatorServer
}

func NewRPCServer(logger rotator.Logger, app rotator.App, host, port string) *Server {
	s := grpc.NewServer(
		grpc.UnaryInterceptor(
			grpc_middleware.ChainUnaryServer(
				grpc_zap.UnaryServerInterceptor(logger.Lgr()),
			),
		),
	)

	internalGrpc := &Server{app: app, logger: logger, srv: s, endpoint: net.JoinHostPort(host, port)}
	gw.RegisterBannersRotatorServer(s, internalGrpc)

	return internalGrpc
}

func (s *Server) Start() error {
	l, err := net.Listen("tcp", s.endpoint)
	if err != nil {
		return fmt.Errorf("listen grpc endpoint -> %w", err)
	}

	if err = s.srv.Serve(l); err != nil {
		return fmt.Errorf("start serve grpc -> %w", err)
	}

	return nil
}

func (s *Server) Stop() {
	s.srv.GracefulStop()
}

func (s *Server) CreateSlot(ctx context.Context, in *gw.Slot) (*gw.Slot, error) {
	if in.Description == "" {
		return nil, status.Errorf(codes.InvalidArgument, "%s: incorrect description", ErrBadRequest)
	}

	slot, err := s.app.CreateSlot(in.Description)
	if err != nil {
		s.logger.Error(fmt.Sprintf("create slot handler -> %s", err))

		return nil, status.Errorf(codes.Internal, "internal server error")
	}

	return &gw.Slot{Id: slot.ID, Description: slot.Description}, nil
}

func (s *Server) CreateBanner(ctx context.Context, in *gw.Banner) (*gw.Banner, error) {
	if in.Description == "" {
		return nil, status.Errorf(codes.InvalidArgument, "%s: incorrect description", ErrBadRequest)
	}

	banner, err := s.app.CreateBanner(in.Description)
	if err != nil {
		s.logger.Error(fmt.Sprintf("create banner handler -> %s", err))

		return nil, status.Errorf(codes.Internal, "internal server error")
	}

	return &gw.Banner{Id: banner.ID, Description: banner.Description}, nil
}

func (s *Server) CreateGroup(ctx context.Context, in *gw.Group) (*gw.Group, error) {
	if in.Description == "" {
		return nil, status.Errorf(codes.InvalidArgument, "%s: incorrect description", ErrBadRequest)
	}

	group, err := s.app.CreateGroup(in.Description)
	if err != nil {
		s.logger.Error(fmt.Sprintf("create gorup handler -> %s", err))

		return nil, status.Errorf(codes.Internal, "internal server error")
	}

	return &gw.Group{Id: group.ID, Description: group.Description}, nil
}

func (s *Server) CreateRotation(ctx context.Context, in *gw.Rotation) (*gw.Message, error) {
	if in.SlotId <= 0 {
		return nil, status.Errorf(codes.InvalidArgument, "%s: incorrect slot id", ErrBadRequest)
	}

	if in.BannerId <= 0 {
		return nil, status.Errorf(codes.InvalidArgument, "%s: incorrect banner id", ErrBadRequest)
	}

	err := s.app.CreateRotation(in.SlotId, in.BannerId)
	if err != nil {
		s.logger.Error(fmt.Sprintf("create rotation handler -> %s", err))

		return nil, status.Errorf(codes.Internal, err.Error())
	}

	return &gw.Message{Message: "Rotation was created"}, nil
}

func (s *Server) DeleteRotation(ctx context.Context, in *gw.Rotation) (*gw.Message, error) {
	if in.SlotId <= 0 {
		return nil, status.Errorf(codes.InvalidArgument, "%s: incorrect slot id", ErrBadRequest)
	}

	if in.BannerId <= 0 {
		return nil, status.Errorf(codes.InvalidArgument, "%s: incorrect banner id", ErrBadRequest)
	}

	err := s.app.DeleteRotation(in.SlotId, in.BannerId)
	if err != nil {
		s.logger.Error(fmt.Sprintf("delete rotation handler -> %s", err))

		return nil, status.Errorf(codes.Internal, err.Error())
	}

	return &gw.Message{Message: "Rotation was deleted"}, nil
}

func (s *Server) CreateClickEvent(ctx context.Context, in *gw.ClickEvent) (*gw.Message, error) {
	if in.SlotId <= 0 {
		return nil, status.Errorf(codes.InvalidArgument, "%s: incorrect slot id", ErrBadRequest)
	}

	if in.BannerId <= 0 {
		return nil, status.Errorf(codes.InvalidArgument, "%s: incorrect banner id", ErrBadRequest)
	}

	if in.GroupId <= 0 {
		return nil, status.Errorf(codes.InvalidArgument, "%s: incorrect group id", ErrBadRequest)
	}

	err := s.app.CreateClickEvent(in.SlotId, in.BannerId, in.GroupId)
	if err != nil {
		s.logger.Error(fmt.Sprintf("create click event handler -> %s", err))

		return nil, status.Errorf(codes.Internal, err.Error())
	}

	return &gw.Message{Message: "Click event was registered"}, nil
}

func (s *Server) BannerForSlot(ctx context.Context, in *gw.SlotRequest) (*gw.Banner, error) {
	if in.SlotId <= 0 {
		return nil, status.Errorf(codes.InvalidArgument, "%s: incorrect slot id", ErrBadRequest)
	}

	if in.GroupId <= 0 {
		return nil, status.Errorf(codes.InvalidArgument, "%s: incorrect group id", ErrBadRequest)
	}

	banner, err := s.app.BannerForSlot(in.SlotId, in.GroupId)
	if err != nil {
		s.logger.Error(fmt.Sprintf("get banner for slot handler -> %s", err))

		return nil, status.Errorf(codes.Internal, err.Error())
	}

	return &gw.Banner{Id: banner.ID, Description: banner.Description}, nil
}
