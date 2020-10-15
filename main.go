package main

import (
	"fmt"
	"github.com/shijting/web/inits"
	"log"
)

const configPath ="configs/config.yaml"

func main()  {
	err := inits.Init(configPath)
	if err !=nil {
		log.Fatal(err)
	}
	fmt.Println(inits.Conf.RedisConfig.DB)
	fmt.Println(inits.Conf.RedisConfig.PoolSize)
}
