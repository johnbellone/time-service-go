package time_server_v1

import (
	"context"

	pb "github.com/johnbellone/time-service/gen/time/v1"

	"github.com/opentracing/opentracing-go"

	"github.com/opentracing/opentracing-go"
	"github.com/grpc-ecosystem/go-grpc-middleware/logging/zap/ctxzap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/codes"
)

type Server struct {
	pb.UnimplementedTimeServer
}

func (s *Server) GetTime(ctx context.Context, req *pb.TimeRequest) (*pb.TimeResponse, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "TimeServer/GetTime")
	defer span.Finish()

	l := ctxzap.Extract(ctx)

	_, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, grpc.Errorf(codes.InvalidArgument, "missing metadata")
	}

	return nil, grpc.Errorf(codes.Unimplemented, "method not implemented")
}
