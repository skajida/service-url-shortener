package main

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"sync"
	"url-shortener/internal/controller"
	"url-shortener/internal/entity"
	"url-shortener/internal/interceptor"
	"url-shortener/internal/pb"
	"url-shortener/internal/service"

	"google.golang.org/grpc/reflection"

	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/recovery"
	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func runGrpc(logger *zap.Logger, appCfg appConfig, controller pb.UrlShortenerServer) {
	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", appCfg.GrpcPort))
	if err != nil {
		logger.Fatal("gRPC listener initialization", zap.Error(err))
	}

	grpcServer := grpc.NewServer(
		grpc.UnaryInterceptor(
			recovery.UnaryServerInterceptor(interceptor.RecoveryOpts(logger)),
		),
	)
	pb.RegisterUrlShortenerServer(grpcServer, controller)
	reflection.Register(grpcServer)

	logger.Info("gRPC server is listening", zap.Uint16("port", appCfg.GrpcPort))
	if err = grpcServer.Serve(listener); err != nil {
		logger.Fatal("Internal gRPC", zap.Error(err))
	}
	grpcServer.GracefulStop()
}

func runRest(ctx context.Context, logger *zap.Logger, appCfg appConfig) {
	mux := runtime.NewServeMux()
	opts := []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())}
	err := pb.RegisterUrlShortenerHandlerFromEndpoint(ctx, mux, fmt.Sprintf("localhost:%d", appCfg.GrpcPort), opts)
	if err != nil {
		logger.Fatal("gRPC client initialization")
	}

	logger.Info("REST server is listening", zap.Uint16("port", appCfg.Port))
	if err := http.ListenAndServe(fmt.Sprintf(":%d", appCfg.Port), mux); err != nil {
		logger.Fatal("Internal REST", zap.Error(err))
	}
}

func main() {
	logger := zap.Must(zap.NewProduction())
	appCfg := initConfig(logger)

	var wg sync.WaitGroup
	defer wg.Wait()
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	gen := entity.NewTokenGenerator(ctx, &wg)

	repository := initRepository(logger, appCfg)
	urlShortener := service.NewUrlShortener(gen, repository)
	server := controller.NewServer(logger, urlShortener)

	wg.Add(1)
	go func(wg *sync.WaitGroup) {
		defer wg.Done()
		runRest(ctx, logger, appCfg)
	}(&wg)
	runGrpc(logger, appCfg, server)
}
