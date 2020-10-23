package middlewares

import (
	"context"
	"github.com/shijting/web/inits/logger"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type validator interface {
	Validate() error
}


func GrpcValidate(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler, ) (resp interface{}, err error) {
	if v, ok := req.(validator); ok {
		if validateErr := v.Validate(); validateErr != nil {
			logger.GetLogger().
				WithField("method", info.FullMethod).
				WithField("query", req).
				Warn(validateErr.Error())
			err = status.Errorf(codes.InvalidArgument, "invalid:%v", validateErr.Error())
			return
		}
	}
	return handler(ctx, req)
}
