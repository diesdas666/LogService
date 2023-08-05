package repo

import (
	"context"
	"errors"
	"example_consumer/internal/core/app"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.uber.org/zap"
)

type AddrBookRepo struct {
	coll *mongo.Collection
}

func NewAddrBookRepo(db *mongo.Database) *AddrBookRepo {
	coll := db.Collection("contacts")
	return &AddrBookRepo{
		coll: coll,
	}
}

type ContactWithPhonesEntity struct {
	ID        primitive.ObjectID `bson:"_id, omitempty"`
	FirstName string             `bson:"firstName"`
	LastName  string             `bson:"lastName"`
	Phones    []*PhoneEntity     `bson:"phones"`
}

type PhoneEntity struct {
	PhoneType   string `bson:"type"`
	PhoneNumber string `bson:"number"`
}

func (r *AddrBookRepo) AddContact(ctx context.Context, c *ContactWithPhonesEntity) (*ContactWithPhonesEntity, error) {
	if c.ID == primitive.NilObjectID {
		v := *c
		v.ID = primitive.NewObjectID()
		c = &v
	}
	_, err := r.coll.InsertOne(ctx, c)
	if err != nil {
		err = fmt.Errorf("error inserting contact into collection: %w", err)
		zap.S().Errorln(err)
		return nil, err
	}
	return c, nil
}

func (r *AddrBookRepo) UpdateContact(ctx context.Context, c *ContactWithPhonesEntity) (found bool, err error) {
	filter := bson.D{{"_id", c.ID}}
	result, err := r.coll.ReplaceOne(ctx, filter, c)
	if err != nil {
		return false, err
	}
	if result.ModifiedCount == 0 {
		return false, nil
	}
	return true, nil
}

func (r *AddrBookRepo) SelectContactByID(ctx context.Context, ID primitive.ObjectID) (*ContactWithPhonesEntity, error) {
	filter := bson.M{"_id": ID}
	var c ContactWithPhonesEntity
	err := r.coll.FindOne(ctx, filter).Decode(&c)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, nil
		}
		err = fmt.Errorf("error fetching contact by id: %v", ID.Hex())
		zap.S().Errorln(err)
		return nil, err
	}
	return &c, nil
}

func (r *AddrBookRepo) SelectAllContacts(ctx context.Context) ([]*ContactWithPhonesEntity, error) {
	// Set options for the find operation
	findOptions := options.Find()

	// Find all documents in the collection
	cursor, err := r.coll.Find(ctx, bson.M{}, findOptions)
	if err != nil {
		return nil, err
	}

	// Iterate over the cursor and decode each document into a bson.M type
	var results []*ContactWithPhonesEntity
	for cursor.Next(context.Background()) {
		var result ContactWithPhonesEntity
		err := cursor.Decode(&result)
		if err != nil {
			app.Logger(ctx).Errorln("failed to decode contact record:", err)
			return nil, err
		}
		results = append(results, &result)
	}

	// Close the cursor to free resources
	err = cursor.Close(context.Background())
	if err != nil {
		return nil, err
	}

	// Return the results
	return results, nil
}

func (r *AddrBookRepo) DeleteContact(ctx context.Context, ID primitive.ObjectID) (found bool, err error) {
	filter := bson.M{"_id": ID}
	result, err := r.coll.DeleteOne(context.Background(), filter)
	println(result)
	if err != nil {
		app.Logger(ctx).Errorln("failed to delete contact record:", err)
		return false, err
	}
	return result.DeletedCount > 0, nil
}
