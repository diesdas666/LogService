package usecase

import (
	"example_consumer/internal/core/outport"
)

type UseCases struct {
	AddrBook outport.AddrBook
	// other output/secondary ports can be added here
}
