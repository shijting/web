package main

import (
	"context"
	"fmt"
	"github.com/golang/glog"
	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/rakyll/statik/fs"
	"github.com/shijting/web/inits"
	"github.com/shijting/web/middlewares"
	"github.com/shijting/web/protos"
	"github.com/shijting/web/services/users"
	_ "github.com/shijting/web/statik"
	"google.golang.org/grpc"
	"log"
	"net"
	"net/http"
	"strings"
)
// config file path
const configPath ="configs/config.yaml"

func main()  {
	err := inits.InitConfig(configPath)
	inits.InitLogger()
	if err !=nil {
		glog.Fatal(err)
	}
	quit := make(chan error, 0)
	// 注册grpc服务
	go runGrpcServer(quit)
	// grpc-gateway服务
	go runGrpcGwServer(quit)

	err = <- quit
	glog.Error(err.Error())
	inits.GetLogger().Fatal(err.Error())
}
func runGrpcServer(quit chan error)  {
	port := inits.Conf.GrpcServerConfig.Port
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
	protos.RegisterUserServiceServer(serv, new(users.UserServiceImpl))
	log.Println("Serving gRPC on port:",port)
	err = serv.Serve(lis)
	if err !=nil {
		quit <- err
	}
}
func runGrpcGwServer(quit chan error) {
	mux := runtime.NewServeMux()
	opts := []grpc.DialOption{grpc.WithInsecure()}
	grpcPort := inits.Conf.GrpcServerConfig.Port
	conn, err := grpc.DialContext(
		context.Background(),
		fmt.Sprintf(":%d", grpcPort),
		opts...,
	)
	if err != nil {
		quit <- fmt.Errorf("failed to dial server: %w", err)
	}
	err = protos.RegisterUserServiceHandler(context.Background(), mux, conn)
	if err != nil {
		quit <- err
	}

	oa := getOpenAPIHandler(quit)

	grpcGwPort := inits.Conf.GrpcGwServerConfig.Port
	if grpcGwPort == 0 {
		grpcGwPort = 8001
	}

	gwServer := &http.Server{
		Addr: fmt.Sprintf(":%d", grpcGwPort),
		Handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// 所有的对外接口都以/api开头
			if strings.HasPrefix(r.URL.Path, "/api") {
				mux.ServeHTTP(w, r)
				return
			}
			// swagger OpenAPI文档服务
			oa.ServeHTTP(w, r)
		}),
	}
	log.Println("Serving gRPC-gateway on port:",grpcGwPort)
	err = gwServer.ListenAndServe()
	if err !=nil {
		quit <- err
	}
}

func getOpenAPIHandler(quit chan error) http.Handler {
	statikFS, err := fs.New()
	if err != nil {
		quit <- err
	}
	// Serve the contents over HTTP.
	return http.FileServer(statikFS)
}
