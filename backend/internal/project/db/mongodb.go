package db

import (
	"backend/internal/project"
	"backend/pkg/logging"
	"context"
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

func (d db) Create(ctx context.Context, project project.Project) (string, error) {
	//TODO implement me
	panic("implement me")
}

func (d db) FindOne(ctx context.Context, id string) (project.Project, error) {
	//TODO implement me
	panic("implement me")
}

func (d db) FindAll(ctx context.Context) ([]project.Project, error) {
	//TODO implement me
	panic("implement me")
}

func (d db) FindByField(ctx context.Context, field, value string) (project.Project, error) {
	//TODO implement me
	panic("implement me")
}

func (d db) Update(ctx context.Context, project project.Project) error {
	//TODO implement me
	panic("implement me")
}

func (d db) Delete(ctx context.Context, id string) error {
	//TODO implement me
	panic("implement me")
}
