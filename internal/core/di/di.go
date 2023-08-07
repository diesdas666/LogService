package di

import (
	"example_consumer/internal/core/app"
	"example_consumer/internal/core/usecase"
)

type DI struct {
	Close    func()
	Config   *app.Config
	UseCases *usecase.UseCases
}
