package time_server_v2

import (
	"context"
	"time"

	pb "github.com/johnbellone/time-service-go/gen/time/v2beta1"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/opentracing/opentracing-go"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
)

type Server struct {
	*pb.UnimplementedTimeServer
}

func NewServer() *Server {
	return &Server{}
}

func (s *Server) GetCurrentTime(ctx context.Context, req *pb.GetTimeRequest) (*pb.GetTimeResponse, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "TimeServer/GetCurrentTime")
	defer span.Finish()

	_, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, grpc.Errorf(codes.InvalidArgument, "missing metadata")
	}

	now := time.Now().UTC()
	rsp := pb.GetTimeResponse{
		Timestamp: timestamppb.New(now),
	}
	return &rsp, nil
}

func (s *Server) GetLocalTime(ctx context.Context, req *pb.GetTimeRequest) (*pb.GetTimeResponse, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "TimeServer/GetLocalTime")
	defer span.Finish()

	_, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, grpc.Errorf(codes.InvalidArgument, "missing metadata")
	}

	rsp := pb.GetTimeResponse{
		Timestamp: timestamppb.New(time.Now()),
	}
	return &rsp, nil
}
