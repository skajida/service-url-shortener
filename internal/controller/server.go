package controller

import (
	"context"
	"url-shortener/internal/model"
	"url-shortener/internal/pb"

	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Server struct {
	logger *zap.Logger
	pb.UnsafeUrlShortenerServer
	service urlShortener
}

func NewServer(logger *zap.Logger, service urlShortener) *Server {
	return &Server{logger: logger, service: service}
}

func (s Server) ReduceUrl(ctx context.Context, url *pb.OriginUrl) (*pb.ShortUrl, error) {
	shortUrl, err := s.service.Generate(ctx, url.GetOriginUrl())
	switch err {
	case nil:
		return &pb.ShortUrl{ShortUrl: shortUrl}, nil
	case model.ErrOriginConflict:
		s.logger.Error("Reducing", zap.Error(err))
		return &pb.ShortUrl{}, status.Error(codes.AlreadyExists, "origin url entry already exists")
	case model.ErrUrlGeneratorInternal:
		fallthrough
	default:
		s.logger.Error("Internal reducing", zap.Error(err))
		return &pb.ShortUrl{}, status.Error(codes.Internal, "internal service error")
	}
}

func (s Server) GetOriginUrl(ctx context.Context, shortUrl *pb.ShortUrl) (*pb.OriginUrl, error) {
	originUrl, err := s.service.LookUp(ctx, shortUrl.GetShortUrl())
	if err != nil {
		s.logger.Error("Getting origin", zap.Error(err))
		return &pb.OriginUrl{}, status.Error(codes.InvalidArgument, "shorted url entry doesn't exist")
	}
	return &pb.OriginUrl{OriginUrl: originUrl}, nil
}
