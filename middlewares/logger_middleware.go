package middlewares

import (
	"context"
	"github.com/shijting/web/inits/logger"
	"google.golang.org/grpc"
	"time"
)

func GrpcLogger(
	ctx context.Context, req interface{},
	info *grpc.UnaryServerInfo,
	handler grpc.UnaryHandler,
) (resp interface{}, err error) {
	timeStart := time.Now()
	resp, err =handler(ctx, req)
	timeElapsed := time.Now().Sub(timeStart)
	logger.GetLogger().
		WithField("method", info.FullMethod).
		WithField("query", req).
		WithField("duration", timeElapsed).
		Info()
	return
}
