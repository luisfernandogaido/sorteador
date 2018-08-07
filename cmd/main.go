package main

import (
	"bitbucket.org/luisfernandogaido/sorteador/app"
	"github.com/gomodule/redigo/redis"
	"log"
)

func main() {
	rd, err := redis.Dial("tcp", "localhost:6379")
	if err != nil {
		log.Fatal(err)
	}
	app.SessaoRedis(rd)
	if err := app.Start(":4003"); err != nil {
		log.Fatal(err)
	}
}
