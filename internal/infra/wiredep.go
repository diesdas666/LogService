package infra

import (
	"example_consumer/internal/core/app"
	"example_consumer/internal/core/di"
	"example_consumer/internal/core/usecase"
	"go.uber.org/zap"
)

func wireDependencies(cfg *app.Config) *di.DI {
	zap.S().Info("Initialize DI objects")
	newDI := &di.DI{
		Config:   cfg,
		UseCases: &usecase.UseCases{},
	}

	cache, cacheCleanup := wireCachePorts(cfg, newDI)

	persistCleanup := wirePersistPorts(
		cfg,
		cache,
		newDI,
	)

	newDI.Close = func() {
		zap.S().Info("Performing cleanup of all initialized DI objects")
		persistCleanup()
		cacheCleanup()
	}
	return newDI
}
