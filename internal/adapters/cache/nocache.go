package cache

import (
	"context"
	"example_consumer/internal/core/outport"
	"fmt"
)

type noCacheAdapter struct {
	chunks map[string]struct{}
}

func NewNoCache() outport.Cache {
	return &noCacheAdapter{
		chunks: make(map[string]struct{}),
	}
}

func (adp *noCacheAdapter) Close() {
	// Nothing to do
}

func (adp *noCacheAdapter) Set(_ context.Context, key outport.CacheKey, _ any) {
	adp.mustHaveCacheChunk(key.Namespace)
}

func (adp *noCacheAdapter) Get(_ context.Context, key outport.CacheKey, _ any) bool {
	adp.mustHaveCacheChunk(key.Namespace)
	return false
}

func (adp *noCacheAdapter) Del(_ context.Context, key outport.CacheKey) {
	adp.mustHaveCacheChunk(key.Namespace)
}

func (adp *noCacheAdapter) Register(partition *outport.CachePartition) {
	ns := partition.Namespace
	if _, ok := adp.chunks[ns]; ok {
		panic(fmt.Sprintf("cache partition with namespace=%s was already registered", ns))
	}
	adp.chunks[ns] = struct{}{}
}

func (adp *noCacheAdapter) mustHaveCacheChunk(ns string) {
	if _, ok := adp.chunks[ns]; !ok {
		panic(fmt.Sprintf("cache partition with namespace=%s was not registered", ns))
	}
}
