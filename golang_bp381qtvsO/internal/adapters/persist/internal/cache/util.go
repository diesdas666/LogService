package cache

import (
	"context"
	"crypto/sha1"
	"encoding/hex"
	"example_consumer/internal/core/outport"
	"fmt"
	"io"
)

func Get[T any](ctx context.Context, cache outport.Cache, key outport.CacheKey) *T {
	value := new(T)
	if cache.Get(ctx, key, &value) {
		return value
	}
	return nil
}

func Set[T any](ctx context.Context, cache outport.Cache, key outport.CacheKey, value *T) {
	cache.Set(ctx, key, value)
}

// GetOrBuild fetches value from cache first and returns it if cache has it already
// Otherwise builder callback function will be invoked
func GetOrBuild[T any](
	ctx context.Context, cache outport.Cache, key outport.CacheKey, builder func() (value *T, doNotSave bool, err error),
) (*T, error) {
	value := new(T)
	if cache.Get(ctx, key, value) {
		return value, nil
	}
	newValue, doNotSave, err := builder()
	if err != nil {
		return nil, err
	}
	if !doNotSave {
		cache.Set(ctx, key, newValue)
	}
	return newValue, nil
}

// BuildCacheKey creates a key that contains two parts - namespace and key itself.
// If namespace+key is shorter than 100 characters then it will be stored in a basic format such as "namespace:key",
// otherwise the key part will be encoded in SHA1. This is due to the fact that usually long keys are
// reducing cache performance (e.g. this is the case with Redis).
func BuildCacheKey(namespace string, key string) outport.CacheKey {
	var newKey string
	if len(namespace)+len(key) < 100 {
		newKey = key
	} else {
		h := sha1.New()
		_, _ = io.WriteString(h, key)
		newKey = hex.EncodeToString(h.Sum(nil))
	}
	return outport.CacheKey{
		Namespace:  namespace,
		EncodedKey: fmt.Sprintf("%s:%s", namespace, newKey),
	}
}
