package db

import (
	"backend/internal/apperror"
	"backend/internal/auth"
	"backend/pkg/logging"
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type db struct {
	collection *mongo.Collection
	logger     *logging.Logger
}

func NewStorage(database *mongo.Database, collection string, logger *logging.Logger) auth.Storage {
	return &db{
		collection: database.Collection(collection),
		logger:     logger,
	}
}

func (d *db) SingIn(ctx context.Context, auth auth.Tokens) error {
	_, err := d.collection.InsertOne(ctx, auth)
	if err != nil {
		return err
	}
	return nil
}

func (d *db) FindByField(ctx context.Context, field, value string) (authData auth.Respond, err error) {
	result := d.collection.FindOne(ctx, bson.M{field: value})
	if result.Err() != nil {
		return authData, nil
	}
	if err := result.Decode(&authData); err != nil {
		return authData, nil
	}
	return authData, nil
}

func (d *db) Update(ctx context.Context, auth auth.Tokens) error {
	oid, err := primitive.ObjectIDFromHex(auth.ID)
	if err != nil {
		return fmt.Errorf("failed to convert token ID to objectID, id: %s", auth.ID)
	}
	filter := bson.M{"_id": oid}
	tokensBytes, err := bson.Marshal(auth)
	if err != nil {
		return fmt.Errorf("failed to marshal token. error: %v", err)
	}
	var updateTokensObj bson.M
	err = bson.Unmarshal(tokensBytes, &updateTokensObj)
	if err != nil {
		return fmt.Errorf("failed unmarshal to tokens bytes. error: %v", err)
	}
	delete(updateTokensObj, "_id")
	update := bson.M{
		"$set": updateTokensObj,
	}
	result, err := d.collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return fmt.Errorf("failed to execute update tokens query. error: %v", err)
	}
	if result.MatchedCount == 0 {
		return apperror.ErrNotFound
	}
	d.logger.Tracef("Matched %d documents and Modified %d documents", result.MatchedCount, result.ModifiedCount)
	return nil
}

func (d *db) Delete(ctx context.Context, field, value string) error {
	filter := bson.M{field: value}
	result, err := d.collection.DeleteOne(ctx, filter)
	if err != nil {
		return fmt.Errorf("failed to execute query. error: %v", err)
	}
	if result.DeletedCount == 0 {
		return apperror.ErrNotFound
	}
	d.logger.Tracef("Deleted %d documents", result.DeletedCount)
	return nil
}
