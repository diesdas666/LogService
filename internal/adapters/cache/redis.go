package cache

import (
	"context"
	"example_consumer/internal/core/app"
	"example_consumer/internal/core/outport"
	"fmt"
	"github.com/redis/go-redis/v9"
	"github.com/vmihailenco/msgpack/v5"
	"go.uber.org/zap"
	"time"
)

type redisCacheAdapter struct {
	client *redis.Client
	chunks map[string]redisCacheChunk
}

type redisCacheChunk struct {
	partition *outport.CachePartition
}

// this context is used to launch redis operations in goroutines without cancellations that http
// request may have. feel free to copy this implementation into some utility go file if you find
// it useful in other parts of the code and don't forget to remote it once this ticket is done:
// https://github.com/golang/go/issues/40221
type contextWithoutDeadline struct {
	ctx context.Context
}

func (*contextWithoutDeadline) Deadline() (time.Time, bool) { return time.Time{}, false }
func (*contextWithoutDeadline) Done() <-chan struct{}       { return nil }
func (*contextWithoutDeadline) Err() error                  { return nil }

func (l *contextWithoutDeadline) Value(key any) any {
	return l.ctx.Value(key)
}

func NewRedisCache(cfg *app.RedisConfig) outport.Cache {
	zap.S().Infof("Open redis client for cache addr=%s db=%d", cfg.Addr, cfg.DB)
	client := redis.NewClient(&redis.Options{
		Addr: cfg.Addr,
		DB:   cfg.DB,
	})
	status := client.Ping(context.Background())
	if status.Err() != nil {
		zap.S().Fatalln("failed to connect to redis:", status.Err())
	}
	return &redisCacheAdapter{
		client: client,
		chunks: make(map[string]redisCacheChunk),
	}
}

func (adp *redisCacheAdapter) Close() {
	if adp.client != nil {
		zap.S().Debug("Close redis client")
		err := adp.client.Close()
		if err != nil {
			zap.S().Warn("Failed to properly close redis cache:", err)
		}
	}
}

func (adp *redisCacheAdapter) Register(partition *outport.CachePartition) {
	ns := partition.Namespace
	zap.S().Debug("Register partition in redis cache:", partition.String())
	if _, ok := adp.chunks[ns]; ok {
		zap.S().Panicf("Cache partition with namespace=%s was already registered", ns)
	}
	adp.chunks[ns] = redisCacheChunk{
		partition: partition,
	}
}

func (adp *redisCacheAdapter) mustGetCacheChunk(ns string) redisCacheChunk {
	if chunk, ok := adp.chunks[ns]; ok {
		return chunk
	}
	panic(fmt.Sprintf("Cache partition with namespace=%s was not registered", ns))
}

func (adp *redisCacheAdapter) Set(ctx context.Context, key outport.CacheKey, value any) {
	zap.S().Debugf("Set item in redis by key=%s", key.String())
	chunk := adp.mustGetCacheChunk(key.Namespace)
	encodedKey := key.EncodedKey
	valueData, err := msgpack.Marshal(value)
	if err != nil {
		app.Logger(ctx).Errorf("Serialize value in redis by key=%s failed with error=%s", key, err)
		return
	}
	runctx := &contextWithoutDeadline{ctx}
	go func() {
		cmd := adp.client.Set(runctx, encodedKey, valueData, chunk.partition.Ttl)
		if err = cmd.Err(); err != nil {
			app.Logger(runctx).Errorf("Async set record to redis by key=%s failed with error: %v", key, err)
		} else {
			app.Logger(runctx).Debugf("Async successfully set data to redis by key=%s", key)
		}
	}()
}

func (adp *redisCacheAdapter) Get(ctx context.Context, key outport.CacheKey, value any) bool {
	app.Logger(ctx).Debugf("Get item in redis by key=%s", key)
	cmd := adp.client.Get(ctx, key.EncodedKey)
	valueData, err := cmd.Bytes()
	if err != nil {
		if err != redis.Nil {
			app.Logger(ctx).Errorf("Read redis data by key=%s failed with error: %v", key, err)
		} else {
			app.Logger(ctx).Debugf("Redis cache miss by key=%s", key)
		}
		return false
	} else {
		app.Logger(ctx).Debugf("Redis cache hit by key=%s", key)
	}
	err = msgpack.Unmarshal(valueData, value)
	if err != nil {
		app.Logger(ctx).Errorf("Unserialize data in redis by key=%s failed with error: %v", key, err)
		return false
	}
	return true
}

func (adp *redisCacheAdapter) Del(ctx context.Context, key outport.CacheKey) {
	encodedKey := key.EncodedKey
	app.Logger(ctx).Debugf("Delete item in redis by key=%s", key)
	runctx := &contextWithoutDeadline{ctx}
	go func() {
		cmd := adp.client.Del(runctx, encodedKey)
		if err := cmd.Err(); err != nil {
			app.Logger(ctx).Errorf("Delete item in redis by key=%s failed with error: %v", key, err)
		} else {
			app.Logger(runctx).Debugf("Async successfully deleted data in redis by key=%s", key)
		}
	}()
}
