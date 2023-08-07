package infra

import (
	"example_consumer/internal/adapters/persist"
	"example_consumer/internal/core/app"
	"example_consumer/internal/core/di"
	"example_consumer/internal/core/outport"
)

func wirePersistPorts(
	cfg *app.Config,
	cache outport.Cache,
	di *di.DI,
) func() {
	pers := persist.NewPersistence(cfg)
	addrBook := persist.NewAddrBookAdapter(
		pers,
		cache,
	)
	di.UseCases.AddrBook = addrBook
	return pers.Close
}
