package main

import (
	"fmt"
	"github.com/golang/glog"
	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	"github.com/shijting/web/gateway"
	"github.com/shijting/web/inits/config"
	"github.com/shijting/web/inits/logger"
	"github.com/shijting/web/inits/psql"
	"github.com/shijting/web/middlewares"
	"github.com/shijting/web/protos"
	"github.com/shijting/web/services/users"
	"google.golang.org/grpc"
	"log"
	"net"
)

const configPath = "configs/config.yaml"

func main() {
	// init config
	if err := config.Init(configPath); err != nil {
		glog.Fatal(err)
	}
	// init logger
	if err := logger.Init(); err != nil {
		glog.Fatal(err)
	}
	// init psql
	if err := psql.Init(); err != nil {
		glog.Fatal(err)
	}
	quit := make(chan error, 0)

	defer func() {
		psql.Close()
		close(quit)
	}()
	// 注册grpc服务
	go runGrpcServer(quit)
	// grpc-gateway服务
	go gateway.RunGrpcGwServer(quit)

	err := <-quit
	glog.Error(err.Error())
	logger.GetLogger().Fatal(err.Error())
}
func runGrpcServer(quit chan error) {
	port := config.Conf.GrpcServerConfig.Port
	if port == 0 {
		port = 8000
	}

	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		quit <- err
	}
	serv := grpc.NewServer(
		grpc.UnaryInterceptor(grpc_middleware.ChainUnaryServer(
			middlewares.GrpcRecover, middlewares.GrpcLogger,middlewares.GrpcValidate,
		)),
	)
	protos.RegisterUserServiceServer(serv, users.NewUserServiceImpl())
	log.Println("Serving gRPC on port:", port)
	err = serv.Serve(lis)
	if err != nil {
		quit <- err
	}
}
