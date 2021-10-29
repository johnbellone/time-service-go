package main

import (
	"flag"
	"fmt"
	"net"
	"os"

	"go.uber.org/zap"

	"github.com/johnbellone/time-service/internal/time-server"

	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	grpc_zap "github.com/grpc-ecosystem/go-grpc-middleware/logging/zap"
	grpc_ctxtags "github.com/grpc-ecosystem/go-grpc-middleware/tags"
	grpc_opentracing "github.com/grpc-ecosystem/go-grpc-middleware/tracing/opentracing"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/health"
	"google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/reflection"
)

const (
	Version = "0.1.0"
)

var (
	Environment string
	GrpcPort    uint
	Verbose     bool
	TlsCertFile string
	TlsKeyFile  string

	BuildTime   string
	GitAbbrv string
	GitCommit   string
)

func init() {
	flag.UintVar(&GrpcPort, "server-port", 9090, "Set the server port.")
	flag.BoolVar(&Verbose, "verbose", false, "Turn on verbose logging.")
	flag.StringVar(&Environment, "environment", "local", "Set the environment name.")
	flag.StringVar(&TlsCertFile, "tls-cert", "server.crt", "Set the path to TLS certificate.")
	flag.StringVar(&TlsKeyFile, "tls-key", "server.key", "Set the path to TLS key.")
}

func main() {
	flag.Parse()

	logger, err := zap.NewProduction()
	if err != nil {
		os.Exit(1)
	}
	defer logger.Sync()

	creds, err := credentials.NewServerTLSFromFile(TlsCertFile, TlsKeyFile)
	if err != nil {
		logger.Fatal(err.Error())
	}

	s := grpc.NewServer(
		grpc.Creds(creds),
		grpc_middleware.WithUnaryServerChain(
			grpc_ctxtags.UnaryServerInterceptor(),
			grpc_opentracing.UnaryServerInterceptor(),
			grpc_zap.UnaryServerInterceptor(logger),
		),
		grpc_middleware.WithStreamServerChain(
			grpc_ctxtags.StreamServerInterceptor(),
			grpc_opentracing.StreamServerInterceptor(),
			grpc_zap.StreamServerInterceptor(logger),
		),
	)

	// Adds all of the handlers for RPC requests to the GRPC server instance. This code is
	// generated when the `protoc` command is run with the `plugins:grpc` switch enabled.
	grpc_health_v1.RegisterHealthServer(s, health.NewServer())
	time_server_v1.RegisterTimeServer(s, time_server_v1.NewServer())
	reflection.Register(s)

	ln, err := net.Listen("tcp", fmt.Sprintf(":%d", Port))
	if err != nil {
		logger.Fatal(err.Error())
	}

	// Set up instance of background context with cancel to gracefully shutdown server when C-c in foreground.
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go func() {
		defer s.GracefulStop()
		<-ctx.Done()
		if err := ln.Close(); err != nil {
			logger.Error(fmt.Sprintf("Failed to close %s: %v", ln.Addr(), err))
		}
	}()

	logger.Info(fmt.Sprintf("%s-%s built at %s", Version, GitAbbrv, BuildTime))
	logger.Info(fmt.Sprintf("server started on %v", ln.Addr()))
	if err = s.Serve(ln); err != nil {
		logger.Fatal(err.Error())
	}
}
