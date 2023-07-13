package utils

import (
	"context"
	"fmt"
	"gin-blog-hufeng/config"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
	"log"
	"time"
)

const REDIS_UTIL_ERR_PREFIX = "utils/redis.go ->"

var (
	ctx = context.Background() // 对服务器的传入请求应该创建一个Context，对服务器的传出调用应该接受一个Context
	// 这个 ctx 只是一个上下文背景，可以存储键值对，在一些context内容可以设置redis的一些限制，但是这里没有
	rdb *redis.Client
)

// Redis 对 Redis 库的操作二次封装， 统一处理错误
//var Redis = new(_redis)

var (
	Redis _redis
)

type _redis struct{}

// InitRedis 初始化 redis 连接
func InitRedis() *redis.Client {
	redisCfg := config.Cfg.Redis
	rdb = redis.NewClient(&redis.Options{
		Addr:     redisCfg.Addr,
		Password: redisCfg.Password,
		DB:       redisCfg.DB,
	})
	// 测试连接状况
	_, err := rdb.Ping(ctx).Result()
	if err != nil {
		log.Panicln("Redis 连接失败：", err)
	}
	log.Println("Redis 连接成功")
	return rdb
}

// GetVal redis 获取值
func (*_redis) GetVal(key string) string {
	return rdb.Get(ctx, key).Val()
}

// Del redis 删除值
func (*_redis) Del(key string) {
	err := rdb.Del(ctx, key).Err()
	if err != nil {
		Logger.Error(REDIS_UTIL_ERR_PREFIX+"Del: ", zap.Error(err))
		panic(err)
	}
}

// SMembers 获取 [集合(Set)] 的成员列表
func (*_redis) SMembers(key string) []string {
	return rdb.SMembers(ctx, key).Val()
}

// Set redis 设置 key value 过期时间
func (*_redis) Set(key string, value interface{}, expiration time.Duration) {
	err := rdb.Set(ctx, key, value, expiration).Err()
	if err != nil {
		Logger.Error(REDIS_UTIL_ERR_PREFIX+"set:", zap.Error(err))
		panic(err)
	}
	fmt.Println("设置redis信息成功")
}

// SAdd 往 [集合(Set)] 中添加 元素
func (*_redis) SAdd(key string, members ...any) {
	err := rdb.SAdd(ctx, key, members...).Err()
	if err != nil {
		Logger.Error(REDIS_UTIL_ERR_PREFIX+"SAdd: ", zap.Error(err))
		panic(err)
	}
}

// GetResult 从 Redis 中取值，不存在会有 redis：nil 的错误
func (*_redis) GetResult(key string) (string, error) {
	return rdb.Get(ctx, key).Result()
}

// SIsMember 判断 元素 是否是 [集合（set）] 的成员
func (*_redis) SIsMember(key string, member any) bool {
	return rdb.SIsMember(ctx, key, member).Val()
}

// HIncrBy 为[哈希表(hash)]中的字段值加上指定增量值(可以为负)
// 如果 key 不存在，自动创建哈希表并执行操作
// 如果 field 不存在，创建该字段值并初始化为 0
func (*_redis) HIncrBy(key, filed string, incr int64) {
	err := rdb.HIncrBy(ctx, key, filed, incr).Err()
	if err != nil {
		Logger.Error(REDIS_UTIL_ERR_PREFIX+"HIncrBy:", zap.Error(err))
		panic(err)
	}
}

// Incr 将 key 中存储的数字 +1
// 如果 key 不存在， 默认初始化为0，再执行 INCR 操作
// 如果 值 包含错误的类型，或是字符串类型的值不能表示为数字，返回错误
func (*_redis) Incr(key string) {
	err := rdb.Incr(ctx, key).Err()
	if err != nil {
		Logger.Error(REDIS_UTIL_ERR_PREFIX+"Incr: ", zap.Error(err))
		panic(err)
	}
}

// GetInt 获取相关内容的数字
func (*_redis) GetInt(key string) int {
	val, _ := rdb.Get(ctx, key).Int() // 文档开始有检查，后面不检查了
	return val
}
