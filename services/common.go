package services

import (
	"context"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc/status"
	"net/http"
)

type ErrorHandler func(ctx context.Context, mux *runtime.ServeMux, marshaler runtime.Marshaler, w http.ResponseWriter, r *http.Request, err error)
func CustomErrorHandler(ctx context.Context, mux *runtime.ServeMux, marshaler runtime.Marshaler, w http.ResponseWriter, r *http.Request, err error)  {
	st := status.Convert(err)
	httpStatus := runtime.HTTPStatusFromCode(st.Code())
	w.WriteHeader(httpStatus)
	w.Write([]byte(st.Message()))
	//mux.errorHandler = Abc()

}
