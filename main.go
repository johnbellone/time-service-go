package main

import (
	"context"
	"flag"
	"fmt"
	"net"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"go.uber.org/zap"

	time_api_v1 "github.com/johnbellone/time-service-go/gen/time/v1"
	time_api_v2 "github.com/johnbellone/time-service-go/gen/time/v2"
	time_server_v1 "github.com/johnbellone/time-service-go/internal/v1"
	time_server_v2 "github.com/johnbellone/time-service-go/internal/v2"

	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	grpc_zap "github.com/grpc-ecosystem/go-grpc-middleware/logging/zap"
	grpc_ctxtags "github.com/grpc-ecosystem/go-grpc-middleware/tags"
	grpc_opentracing "github.com/grpc-ecosystem/go-grpc-middleware/tracing/opentracing"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/insecure"
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

	BuildTime string
	GitAbbrv  string
	GitCommit string
)

func init() {
	flag.UintVar(&GrpcPort, "server-port", 50010, "Set the server port.")
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
	} else if Verbose {
		// TODO: Enable verbose logging settings in the zap top-level zap configuration. 
	}
	defer logger.Sync()

	creds := insecure.NewCredentials()

	// Set the TLS configuration for gRPC server channels. By default, this is an insecure
	// bundle all environments except production. 
	if Environment == "prod" {
		if TlsCertFile == "" || TlsKeyFile == "" {
			logger.Fatal("tls-cert and tls-key must be set in production")
			os.Exit(2)
		}

		if creds, err = credentials.NewServerTLSFromFile(TlsCertFile, TlsKeyFile); err != nil {
			logger.Fatal("failed creating tls credentials", zap.Error(err))
			os.Exit(3)
		}
	}

	// Set up all of the gRPC middleware ahead of time for all endpoints. 
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
	time_api_v1.RegisterTimeServiceServer(s, time_server_v1.NewServer())
	time_api_v2.RegisterTimeServiceServer(s, time_server_v2.NewServer())
	reflection.Register(s)

	ln, err := net.Listen("tcp", fmt.Sprintf(":%d", GrpcPort))
	if err != nil {
		logger.Fatal("failed opening server socket", zap.Error(err))
	}

	// Set up instance of background context with cancel to gracefully shutdown server if
	// C-c in the foreground.
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
	wg := sync.WaitGroup{}

	wg.Add(1)
	go func() {
		defer s.GracefulStop()
		select {
		case <-ctx.Done():
			logger.Info("closing server")
			break
		case s := <-c:
			logger.Info("received signal", zap.Any("signal", s))
			break
		}

		logger.Info("stopping server")

		if err := ln.Close(); err != nil {
			logger.Error("failed closing server socket", zap.Error(err))
		}
		cancel()
		wg.Done()
	}()

	// Print the program's build information at each and every start.
	program, _ := os.Executable()
	logger.Info("build info",
		zap.String("Executable", program),
		zap.String("Version", Version),
		zap.String("GitAbbrv", GitAbbrv),
		zap.String("GitCommit", GitCommit),
		zap.String("BuildTime", BuildTime),
	)

	logger.Info("starting server", zap.Uint("GrpcPort", GrpcPort))
	if err = s.Serve(ln); err != nil {
		logger.Error("server error", zap.Error(err))
	}

	wg.Wait()
	logger.Info("shutdown server")
}
