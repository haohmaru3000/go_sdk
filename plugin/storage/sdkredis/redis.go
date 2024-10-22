package sdkredis

// Go-Redis is an alternative option for Redis
// Github: https://github.com/go-redis/redis
//
// It supports:
// 		Redis 3 commands except QUIT, MONITOR, SLOWLOG and SYNC.
// 		Automatic connection pooling with circuit breaker support.
// 		Pub/Sub.
// 		Transactions.
// 		Pipeline and TxPipeline.
// 		Scripting.
// 		Timeouts.
// 		Redis Sentinel.
// 		Redis Cluster.
// 		Cluster of Redis Servers without using cluster mode and Redis Sentinel.
// 		Ring.
// 		Instrumentation.
// 		Cache friendly.
// 		Rate limiting.
// 		Distributed Locks.

import (
	"context"
	"flag"

	"github.com/haohmaru3000/go_sdk/logger"
	"github.com/redis/go-redis/v9"
)

var (
	defaultRedisName      = "DefaultRedis"
	DefaultRedisDB        = getDefaultRedisDB()
	defaultRedisMaxActive = 0 // 0 is unlimited max active connection
	defaultRedisMaxIdle   = 10
)

type RedisDBOpt struct {
	Prefix    string
	RedisUri  string
	MaxActive int
	MaxIde    int
}

type redisDB struct {
	context context.Context
	name    string
	client  *redis.Client
	logger  logger.Logger
	*RedisDBOpt
}

func getDefaultRedisDB() *redisDB {
	return NewRedisDB(context.Background(), defaultRedisName, "")
}

func NewRedisDB(ctx context.Context, name, flagPrefix string) *redisDB {
	return &redisDB{
		context: ctx,
		name:    name,
		RedisDBOpt: &RedisDBOpt{
			Prefix:    flagPrefix,
			MaxActive: defaultRedisMaxActive,
			MaxIde:    defaultRedisMaxIdle,
		},
	}
}

func (r *redisDB) GetPrefix() string {
	return r.Prefix
}

func (r *redisDB) isDisabled() bool {
	return r.RedisUri == ""
}

func (r *redisDB) InitFlags() {
	prefix := r.Prefix
	if r.Prefix != "" {
		prefix += "-"
	}

	flag.StringVar(&r.RedisUri, prefix+"go-redis-uri", "", "(For go-redis) Redis connection-string. Ex: redis://localhost/0")
	flag.IntVar(&r.MaxActive, prefix+"go-redis-pool-max-active", defaultRedisMaxActive, "(For go-redis) Override redis pool MaxActive")
	flag.IntVar(&r.MaxIde, prefix+"go-redis-pool-max-idle", defaultRedisMaxIdle, "(For go-redis) Override redis pool MaxIdle")
}

func (r *redisDB) Configure() error {
	if r.isDisabled() {
		return nil
	}

	r.logger = logger.GetCurrent().GetLogger(r.name)
	r.logger.Info("Connecting to Redis at ", r.RedisUri, "...")

	opt, err := redis.ParseURL(r.RedisUri)

	if err != nil {
		r.logger.Error("Cannot parse Redis ", err.Error())
		return err
	}

	opt.PoolSize = r.MaxActive
	opt.MinIdleConns = r.MaxIde

	client := redis.NewClient(opt)

	// Ping to test Redis connection
	if err := client.Ping(r.context).Err(); err != nil {
		r.logger.Error("Cannot connect Redis. ", err.Error())
		return err
	}

	// Connect successfully, assign client to goRedisDB
	r.client = client
	return nil
}

func (r *redisDB) Name() string {
	return r.name
}

func (r *redisDB) Get() interface{} {
	return r.client
}

func (r *redisDB) Run() error {
	return r.Configure()
}

func (r *redisDB) Stop() <-chan bool {
	if r.client != nil {
		if err := r.client.Close(); err != nil {
			r.logger.Info("cannot close ", r.name)
		}
	}

	c := make(chan bool)
	go func() { c <- true }()
	return c
}
