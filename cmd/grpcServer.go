package main

import (
	"github.com/shijting/web/inits"
	"github.com/shijting/web/protos"
	"github.com/shijting/web/services/users"
	"google.golang.org/grpc"
	"log"
	"net"
)
// config file path
const configPath ="configs/config.yaml"
func main()  {
	err := inits.InitConfig(configPath)
	if err !=nil {
		log.Fatal(err)
	}
	lis,err := net.Listen("tcp", ":8000")
	if err !=nil {
		log.Fatal(err)
	}
	serv := grpc.NewServer()
	protos.RegisterUserServiceServer(serv, new(users.UsersServiceImpl))

	err = serv.Serve(lis)
	if err !=nil {
		log.Fatal(err)
	}
}
