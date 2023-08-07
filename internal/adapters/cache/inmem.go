package cache

import (
	"context"
	"example_consumer/internal/adapters/cache/internal"
	"example_consumer/internal/core/outport"
	"fmt"
	"github.com/vmihailenco/msgpack/v5"
	"go.uber.org/zap"
)

type inMemCacheAdapter struct {
	chunks map[string]inMemCacheChunk
}

func NewInMemCache() outport.Cache {
	return &inMemCacheAdapter{
		chunks: make(map[string]inMemCacheChunk),
	}
}

type inMemCacheChunk struct {
	partition *outport.CachePartition
	cache     *internal.TinyLFU
}

func (adp *inMemCacheAdapter) Close() {
	// Nothing to do
}

func (adp *inMemCacheAdapter) Register(partition *outport.CachePartition) {
	ns := partition.Namespace
	if _, ok := adp.chunks[ns]; ok {
		panic(fmt.Sprintf("Cache partition with namespace=%s was already registered", ns))
	}
	c := internal.NewTinyLFU(partition.LocalMaxItems, partition.Ttl)
	adp.chunks[ns] = inMemCacheChunk{
		partition: partition,
		cache:     c,
	}
}

func (adp *inMemCacheAdapter) mustGetCacheChunk(ns string) inMemCacheChunk {
	if chunk, ok := adp.chunks[ns]; ok {
		return chunk
	}
	panic(fmt.Sprintf("Cache partition with namespace=%s was not registered", ns))
}

func (adp *inMemCacheAdapter) Set(_ context.Context, key outport.CacheKey, value any) {
	zap.S().Debugf("Set item in in-mem cache by cacheKey=%s", key)
	chunk := adp.mustGetCacheChunk(key.Namespace)
	result, err := msgpack.Marshal(value)
	if err != nil {
		panic(fmt.Errorf("error marshalling value: %w", err))
	}
	chunk.cache.Set(key.EncodedKey, result)
}

func (adp *inMemCacheAdapter) Get(_ context.Context, key outport.CacheKey, value any) bool {
	zap.S().Debugf("Get item from in-mem cache by cacheKey=%s", key)
	chunk := adp.mustGetCacheChunk(key.Namespace)
	if data, ok := chunk.cache.Get(key.EncodedKey); ok {
		err := msgpack.Unmarshal(data, value)
		if err != nil {
			panic(fmt.Errorf("error unmarshalling value: %w", err))
		}
		return true
	}
	return false
}

func (adp *inMemCacheAdapter) Del(_ context.Context, key outport.CacheKey) {
	zap.S().Debugf("Delete item in in-mem cache by cacheKey=%s", key)
	chunk := adp.mustGetCacheChunk(key.Namespace)
	chunk.cache.Del(key.EncodedKey)
}
