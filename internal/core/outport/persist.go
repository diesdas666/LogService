package outport

import (
	"go.mongodb.org/mongo-driver/mongo"
)

type Persistence interface {
	DB() *mongo.Database
	Close()
}
