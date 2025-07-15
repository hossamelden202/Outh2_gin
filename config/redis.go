package config

import (
	"context"
	"fmt"

	"github.com/redis/go-redis/v9"
)

var Rd *redis.Client
var Ctx=context.Background()
func ConnnectRedis(){


	
	Rd=redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
		DB: 0,
		Password: "",
	})
	if err:=Rd.Ping(Ctx).Err();err!=nil{
	panic("connection to redis failed")
	}
// Rd.FlushAll(Ctx)
fmt.Println("connected to redis successfully")
}