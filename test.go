package main

import (
	"fmt"
	"net/http"
)

func main()  {
	_ = http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {

	})
	_ = http.Server{
		Addr:              "",
		Handler:           http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {

		}),
	}
	//http.Handle("/", )
	a := []int{0,1,2,3}
	a = append([]int{-3,-2, -1}, a...)
	fmt.Println(a)
	fmt.Println(append(a[:1], []int{5,6,7}...))
}
