package gateway

import (
	"context"
	"fmt"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/rakyll/statik/fs"
	"github.com/shijting/web/inits/config"
	"github.com/shijting/web/protos"
	_ "github.com/shijting/web/statik"
	"google.golang.org/grpc"
	"google.golang.org/grpc/status"
	"log"
	"mime"
	"net/http"
	"strings"
)

func RunGrpcGwServer(quit chan error) {

	errorOption := runtime.WithErrorHandler(gatewayErrorHandler)
	mux := runtime.NewServeMux(errorOption)

	//mux := runtime.NewServeMux()
	opts := []grpc.DialOption{grpc.WithInsecure()}
	grpcPort := config.Conf.GrpcServerConfig.Port
	ctx := context.Background()
	conn, err := grpc.DialContext(
		ctx,
		fmt.Sprintf(":%d", grpcPort),
		opts...,
	)
	if err != nil {
		quit <- fmt.Errorf("failed to dial server: %w", err)
	}

	err = protos.RegisterUserServiceHandler(ctx, mux, conn)
	if err != nil {
		quit <- err
	}

	oa := getOpenAPIHandler(quit)

	grpcGwPort := config.Conf.GrpcGwServerConfig.Port
	if grpcGwPort == 0 {
		grpcGwPort = 8001
	}
	m := NewMiddleware()
	m.Use(timeMiddleware)
	m.Use(logMiddleware)
	// start timeMiddleware -> start logMiddleware ->content -> end logMiddleware ->end timeMiddleware
	gwServer := &http.Server{
		Addr: fmt.Sprintf(":%d", grpcGwPort),
		Handler: m.middleWareHandler(handler(mux, oa)),
	}
	log.Println("Serving gRPC-gateway on port:",grpcGwPort)
	err = gwServer.ListenAndServe()
	if err !=nil {
		quit <- err
	}
}

func handler(mux *runtime.ServeMux, oa http.Handler ) http.Handler  {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// 所有的对外接口都以/api开头
		if strings.HasPrefix(r.URL.Path, "/api") {
			fmt.Println("content")
			mux.ServeHTTP(w, r)
			return
		}
		// swagger OpenAPI文档服务
		oa.ServeHTTP(w, r)
	})
}
func (r *Middleware)middleWareHandler(h http.Handler) http.Handler  {
	var mergedHandler = h
	for i := len(r.middlewareChain) - 1; i >= 0; i-- {
		mergedHandler = r.middlewareChain[i](mergedHandler)
	}
	return mergedHandler
}

func getOpenAPIHandler(quit chan error) http.Handler {
	err := mime.AddExtensionType(".svg", "image/svg+xml")
	if err != nil {
		quit <- err
	}
	statikFS, err := fs.New()
	if err != nil {
		quit <- err
	}
	return http.FileServer(statikFS)
}

func gatewayErrorHandler(ctx context.Context, mux *runtime.ServeMux, marshaler runtime.Marshaler, w http.ResponseWriter, r *http.Request, err error) {

	st := status.Convert(err)
	httpStatus := runtime.HTTPStatusFromCode(st.Code())
	w.WriteHeader(httpStatus)
	//w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(st.Message()))
}