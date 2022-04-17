package db

import (
	"backend/internal/apperror"
	"backend/internal/project"
	"backend/pkg/logging"
	"context"
	"errors"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type db struct {
	collection *mongo.Collection
	logger     *logging.Logger
}

func NewStorage(database *mongo.Database, collection string, logger *logging.Logger) project.Storage {
	return &db{
		collection: database.Collection(collection),
		logger:     logger,
	}
}

func (d *db) Create(ctx context.Context, project project.Project) (string, error) {
	d.logger.Debug("create project")
	result, err := d.collection.InsertOne(ctx, project)
	if err != nil {
		return "", fmt.Errorf("failed to create project due to error: %v", err)
	}
	d.logger.Debug("convert insertedID to ObjectID")
	oid, ok := result.InsertedID.(primitive.ObjectID)
	if ok {
		return oid.Hex(), nil
	}
	d.logger.Trace(project)
	return "", fmt.Errorf("failed to convert objectID to hex. probbably oid: %v", oid)
}

func (d *db) FindOne(ctx context.Context, id string) (p project.Project, err error) {
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return p, fmt.Errorf("failed to convert hex to objectID, hex: %s", id)
	}
	filter := bson.M{"_id": oid}
	result := d.collection.FindOne(ctx, filter)
	if result.Err() != nil {
		if errors.Is(result.Err(), mongo.ErrNoDocuments) {
			return p, apperror.ErrNotFound
		}
		return p, fmt.Errorf("failed to find one project by id: %s due to error: %v", id, err)
	}
	if err := result.Decode(&p); err != nil {
		return p, fmt.Errorf("failed to decode project (id: %s) from DB due to error: %v", id, err)
	}
	return p, nil
}

func (d *db) FindAll(ctx context.Context) (p []project.Project, err error) {
	cursor, err := d.collection.Find(ctx, bson.M{})
	if cursor.Err() != nil {
		return p, fmt.Errorf("failed to find all projects due to error: %v", err)
	}
	if err := cursor.All(ctx, &p); err != nil {
		return p, fmt.Errorf("failed to read of documents form cursor. error: %s", err)
	}
	return p, nil
}

func (d *db) FindByField(ctx context.Context, field, value string) (p project.Project, err error) {
	result := d.collection.FindOne(ctx, bson.M{field: value})
	if result.Err() != nil {
		return p, nil
	}
	if err := result.Decode(&p); err != nil {
		return p, nil
	}
	return p, nil
}

func (d *db) Update(ctx context.Context, project project.Project) error {
	oid, err := primitive.ObjectIDFromHex(project.ID)
	if err != nil {
		return fmt.Errorf("failed to convert project ID to objectID, id: %s", project.ID)
	}
	filter := bson.M{"_id": oid}
	projectBytes, err := bson.Marshal(project)
	if err != nil {
		return fmt.Errorf("failed to marshal project. error: %v", err)
	}
	var updateProjectObj bson.M
	err = bson.Unmarshal(projectBytes, &updateProjectObj)
	if err != nil {
		return fmt.Errorf("failed unmarshal to project bytes. error: %v", err)
	}
	delete(updateProjectObj, "_id")
	update := bson.M{
		"$set": updateProjectObj,
	}
	result, err := d.collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return fmt.Errorf("failed to execute update project query. error: %v", err)
	}
	if result.MatchedCount == 0 {
		return apperror.ErrNotFound
	}
	d.logger.Tracef("Matched %d documents and Modified %d documents", result.MatchedCount, result.ModifiedCount)
	return nil
}

func (d *db) Delete(ctx context.Context, id string) error {
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return fmt.Errorf("failed to convert project ID to objectID, id: %s", id)
	}
	filter := bson.M{"_id": oid}
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
