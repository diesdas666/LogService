package persist

import (
	"context"
	"example_consumer/internal/core/app"
	"example_consumer/internal/core/outport"
	"fmt"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.uber.org/zap"
)

type dbAdapter struct {
	client *mongo.Client
	db     *mongo.Database
}

// NewPersistence connects to PostgreSQL database and returns Persistence interface that wraps database reference
func NewPersistence(cfg *app.Config) outport.Persistence {
	dbc := cfg.Database
	var connStr string
	if len(dbc.User) > 0 && dbc.User != "_" {
		connStr = fmt.Sprintf("mongodb://%s:%s@%s:%d", dbc.User, dbc.Password, dbc.Host, dbc.Port)
	} else {
		connStr = fmt.Sprintf("mongodb://%s:%d", dbc.Host, dbc.Port)
	}
	zap.S().Infoln("establishing connection to mongodb database...")
	clientOpts := options.Client().ApplyURI(connStr)
	client, err := mongo.Connect(context.TODO(), clientOpts)
	if err != nil {
		zap.S().Fatalln("error connecting to mongodb database:", err)
	}

	err = client.Ping(context.TODO(), nil)
	if err != nil {
		zap.S().Fatalln("failed to ping mongodb database:", err)
	}
	zap.S().Infoln("database initialization was successfully performed")

	db := client.Database(dbc.Name)

	return &dbAdapter{
		client: client,
		db:     db,
	}
}

func (d dbAdapter) DB() *mongo.Database {
	return d.db
}

func (d dbAdapter) Close() {
	if err := d.client.Disconnect(context.TODO()); err != nil {
		zap.S().Fatalln("failed to close mongdb connection: ", err)
	}
}
