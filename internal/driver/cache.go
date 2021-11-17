package driver

// import (
// 	"context"
// 	"fmt"

// 	"github.com/go-redis/redis/v8"
// )

// type CACHE struct {
// 	REDIS *redis.Conn
// }
// var redisConn = &CACHE{}

// func GetRDC(dsn string)(*CACHE, error) {

// 	// val, err := rdb.Get(ctx, "key").Result()
// 	// switch {
// 	// case err == redis.Nil:
// 	// 	fmt.Println("key does not exist")
// 	// case err != nil:
// 	// 	fmt.Println("Get failed", err)
// 	// case val == "":
// 	// 	fmt.Println("value is empty")
// 	// }

// }

// func NewRedis(dsn string) (*redis.Client, error) {
// 	// ctx := context.Background()
// 	rdb := redis.NewClient(&redis.Options{

// 	})

// 	return rdb, nil
// }
