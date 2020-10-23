package gateway

import (
	"fmt"
	"log"
	"net/http"
	"time"
)

type middleware func(next http.Handler) http.Handler

type Middleware struct {
	middlewareChain []middleware
}

func NewMiddleware() *Middleware {
	return &Middleware{}
}

func (r *Middleware) Use(m middleware) {
	r.middlewareChain = append(r.middlewareChain, m)
}

func timeMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(wr http.ResponseWriter, r *http.Request) {
		timeStart := time.Now()
		fmt.Println("start timeMiddleware")
		// next handler
		next.ServeHTTP(wr, r)

		timeElapsed := time.Since(timeStart)
		log.Println(timeElapsed)
		fmt.Println("end timeMiddleware")
	})
}
func logMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(wr http.ResponseWriter, r *http.Request) {
		fmt.Println("start logMiddleware")
		// next handler
		next.ServeHTTP(wr, r)
		fmt.Println("end logMiddleware")
	})
}