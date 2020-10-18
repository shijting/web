package gateway

import (
	"context"
	"fmt"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/rakyll/statik/fs"
	"github.com/shijting/web/inits"
	"github.com/shijting/web/protos"
	_ "github.com/shijting/web/statik"
	"google.golang.org/grpc"
	"log"
	"mime"
	"net/http"
	"strings"
)

func RunGrpcGwServer(quit chan error) {
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
		fmt.Println(111)
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
		fmt.Println(2222)
		quit <- err
	}
}

func getOpenAPIHandler(quit chan error) http.Handler {
	err := mime.AddExtensionType(".svg", "image/svg+xml")
	if err != nil {
		fmt.Println(4444)
		quit <- err
	}
	statikFS, err := fs.New()
	if err != nil {
		fmt.Println(3333)
		quit <- err
	}
	// Serve the contents over HTTP.
	return http.FileServer(statikFS)
}