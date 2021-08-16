package redis

import (
	"chain-demo/config"
	"context"
	"fmt"
	"github.com/go-redis/redis/v8" // 注意导入的是新版本
	"time"
)

// 声明一个全局的rdb变量
var rdb *redis.Client

// 初始化连接
func initClient() (err error) {

	rdb = redis.NewClient(&redis.Options{
		Addr:     conf.Conf.Redis.Server,
		Password: conf.Conf.Redis.Pwd, // no password set
		DB:       0,                   // use default DB
		PoolSize: 100,                 // 连接池大小
	})

	// 连接Redis哨兵模式
	//rdb := redis.NewFailoverClient(&redis.FailoverOptions{
	//	MasterName:    "master",
	//	SentinelAddrs: []string{"x.x.x.x:26379", "xx.xx.xx.xx:26379", "xxx.xxx.xxx.xxx:26379"},
	//})

	// 连接Redis集群
	//rdb := redis.NewClusterClient(&redis.ClusterOptions{
	//	Addrs: []string{":7000", ":7001", ":7002", ":7003", ":7004", ":7005"},
	//})

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err = rdb.Ping(ctx).Result()
	return err
}

// V8Example 测试方法
func V8Example() {
	ctx := context.Background()
	if err := initClient(); err != nil {
		return
	}

	err := rdb.Set(ctx, "key", "value", 0).Err()
	if err != nil {
		panic(err)
	}

	val, err := rdb.Get(ctx, "key").Result()
	if err != nil {
		panic(err)
	}
	fmt.Println("key", val)

	val2, err := rdb.Get(ctx, "key2").Result()
	if err == redis.Nil {
		fmt.Println("key2 does not exist")
	} else if err != nil {
		panic(err)
	} else {
		fmt.Println("key2", val2)
	}
	// Output: key value
	// key2 does not exist
}
