package cache

import (
	"context"
	"example_consumer/internal/core/model"
	"example_consumer/internal/core/outport"
	"time"
)

const nsAddrBookContactByID = "ContactById"

type ContactByIdPartition struct {
	cache outport.Cache
}

// RegisterContactByID registers cache to get/set/delete contact record by ID.
func RegisterContactByID(cache outport.Cache) ContactByIdPartition {
	prt := &outport.CachePartition{
		Namespace:     nsAddrBookContactByID,
		Ttl:           30 * time.Second,
		LocalMaxItems: 1000,
	}
	cache.Register(prt)
	return ContactByIdPartition{cache: cache}
}

// Get returns contact record from cache by ID
func (pr ContactByIdPartition) Get(ctx context.Context, ID string) *model.Contact {
	key := BuildCacheKey(nsAddrBookContactByID, ID)
	return Get[model.Contact](ctx, pr.cache, key)
}

// Set adds/updates contact record by ID
func (pr ContactByIdPartition) Set(ctx context.Context, c *model.Contact) {
	key := BuildCacheKey(nsAddrBookContactByID, c.ID)
	Set(ctx, pr.cache, key, c)
}

// Del deletes contact record by ID
func (pr ContactByIdPartition) Del(ctx context.Context, ID string) {
	key := BuildCacheKey(nsAddrBookContactByID, ID)
	pr.cache.Del(ctx, key)
}
