package middlewares

import (
	"context"
	"github.com/shijting/web/inits"
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
			inits.GetLogger().
				WithField("method", info.FullMethod).
				WithField("query", req).
				Error(err)
			err = status.Errorf(codes.InvalidArgument, "参数无效:%v", r)
		}
	}()
	resp, err =handler(ctx, req)
	return
}
