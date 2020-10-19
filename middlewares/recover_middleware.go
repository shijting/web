package middlewares

import (
	"context"
	"fmt"
	"github.com/shijting/web/inits/logger"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)
// 接管 panic
func GrpcRecover(
	ctx context.Context, req interface{},
	info *grpc.UnaryServerInfo,
	handler grpc.UnaryHandler,
) (resp interface{}, err error) {
	//fmt.Printf("%v---%v\n", info.FullMethod, info.Server)

	defer func() {
		if r := recover(); r != nil {
			// todo
			logger.GetLogger().
				WithField("method", info.FullMethod).
				WithField("query", req).
				Error(err)
			fmt.Println("panic:", r)
			err = status.Errorf(codes.Internal, "系统错误！")
		}
	}()
	resp, err =handler(ctx, req)
	return
}
