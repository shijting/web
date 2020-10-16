package middlewares

import (
	"context"
	"fmt"
	"github.com/shijting/web/inits"
	"google.golang.org/grpc"
)
// 接管panic
func GrpcRecover(
	ctx context.Context, req interface{},
	info *grpc.UnaryServerInfo,
	handler grpc.UnaryHandler,
) (resp interface{}, err error) {
	fmt.Printf("%v---%v\n", info.FullMethod, info.Server)
	defer func() {
		if err := recover(); err != nil {
			// todo
			inits.GetLogger().
				WithField("fullMethod", info.FullMethod).
				WithField("query", req).
				Error(err)
		}
	}()

	return handler(ctx, req)
}
