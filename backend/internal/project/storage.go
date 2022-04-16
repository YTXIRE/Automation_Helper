package project

import "context"

type Storage interface {
	Create(ctx context.Context, project Project) (string, error)
	FindOne(ctx context.Context, id string) (Project, error)
	FindAll(ctx context.Context) ([]Project, error)
	FindByField(ctx context.Context, field, value string) (Project, error)
	Update(ctx context.Context, project Project) error
	Delete(ctx context.Context, id string) error
}
