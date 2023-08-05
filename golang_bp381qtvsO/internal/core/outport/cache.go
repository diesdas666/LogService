package outport

import (
	"context"
	"fmt"
	"time"
)

// Cache declares generic cache interface to be implemented by real cache implementation (such as Redis for example)
type Cache interface {
	Close()
	Register(partition *CachePartition)
	Set(ctx context.Context, key CacheKey, value any)
	Get(ctx context.Context, key CacheKey, value any) bool
	Del(ctx context.Context, key CacheKey)
}

type CachePartition struct {
	Namespace     string        // Namespace key prefix, must be unique amongst other namespaces
	Ttl           time.Duration // Time-to-live
	LocalMaxItems int           // Max number of items in cache (in-memory cache only, for redis etc will be ignored)
}

func (cp *CachePartition) String() string {
	return fmt.Sprintf("{namespace=%s ttl=%s localMaxItems=%d}", cp.Namespace, cp.Ttl.String(), cp.LocalMaxItems)
}

// CacheKey consists of two parts - namespace (can be your entity type) and encoded key itself
// Refer to BuildCacheKey() function is this file for more information.
type CacheKey struct {
	Namespace  string
	EncodedKey string
}

func (k CacheKey) String() string {
	return fmt.Sprintf("{namespace=%s encodedKey=%s}", k.Namespace, k.EncodedKey)
}
