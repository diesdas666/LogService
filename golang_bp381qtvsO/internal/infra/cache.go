package infra

import (
	"example_consumer/internal/adapters/cache"
	"example_consumer/internal/core/app"
	"example_consumer/internal/core/di"
	"example_consumer/internal/core/outport"
	"fmt"
)

func wireCachePorts(cfg *app.Config, _ *di.DI) (outport.Cache, func()) {
	switch cfg.Cache.Type {
	case "none":
		return cache.NewNoCache(), func() {}
	case "inmem":
		return cache.NewInMemCache(), func() {}
	case "redis":
		r := cache.NewRedisCache(&cfg.Cache.Redis)
		return r, func() {
			r.Close()
		}
	default:
		panic(fmt.Sprintf("unknown cache type: %s", cfg.Cache.Type))
	}
}
