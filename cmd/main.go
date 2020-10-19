package main

import (
	"context"
	"fmt"
	"github.com/go-pg/pg/v10"
	"github.com/golang/glog"
	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	"github.com/shijting/web/gateway"
	"github.com/shijting/web/inits/config"
	"github.com/shijting/web/inits/logger"
	"github.com/shijting/web/inits/psql"
	"github.com/shijting/web/middlewares"
	"github.com/shijting/web/protos"
	"github.com/shijting/web/services/users"
	_ "github.com/shijting/web/statik"
	"google.golang.org/grpc"
	"log"
	"net"
)
// config file path
const configPath ="configs/config.yaml"

func main()  {
	// init config
	err := config.Init(configPath)
	if err !=nil {
		glog.Fatal(err)
	}
	// init logger
	err = logger.Init()
	if err !=nil {
		glog.Fatal(err)
	}
	// init psql
	err = psql.Init()
	if err !=nil {
		glog.Fatal(err)
	}
	var version string
	_, err = psql.GetDB().QueryOneContext(context.Background(), pg.Scan(&version), "SELECT version()")
	if err != nil {
		panic(err)
	}
	fmt.Println("version:", version)
	quit := make(chan error, 0)
	// 注册grpc服务
	go runGrpcServer(quit)
	// grpc-gateway服务
	go gateway.RunGrpcGwServer(quit)

	err = <- quit
	glog.Error(err.Error())
	logger.GetLogger().Fatal(err.Error())
}
func runGrpcServer(quit chan error)  {
	port := config.Conf.GrpcServerConfig.Port
	if port == 0 {
		port = 8000
	}

	lis,err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err !=nil {
		quit <- err
	}
	serv := grpc.NewServer(
		grpc.UnaryInterceptor(grpc_middleware.ChainUnaryServer(
			middlewares.GrpcRecover,
		)),
	)
	protos.RegisterUserServiceServer(serv, users.NewUserServiceImpl())
	log.Println("Serving gRPC on port:",port)
	err = serv.Serve(lis)
	if err !=nil {
		quit <- err
	}
}

