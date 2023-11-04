package interceptor

import (
	"context"

	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/recovery"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func RecoveryOpts(logger *zap.Logger) recovery.Option {
	return recovery.WithRecoveryHandlerContext(
		func(ctx context.Context, p any) error {
			logger.Error("Recovery interceptor", zap.Any("%s", p))
			return status.Errorf(codes.Internal, "%s", p)
		},
	)
}
