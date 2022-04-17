package project

import (
	"backend/internal/apperror"
	"backend/pkg/logging"
	"context"
	"fmt"
	"time"
)

type Service struct {
	logger *logging.Logger
}

func (s *Service) Create(ctx context.Context, project Project, storage Storage) (Project, error) {
	err := project.Validate()
	if err != nil {
		return Project{}, apperror.NewAppError(err, err.Error(), "")
	}
	findUser, err := storage.FindByField(ctx, "name", project.Name)
	if err != nil {
		return Project{}, apperror.NewAppError(err, fmt.Sprintf("failed to find project with name: %s. error: %v", project.Name, err), "")
	}
	if findUser.Name != "" {
		return Project{}, apperror.NewAppError(nil, "this name is already busy", "")
	}
	createProjectData := Project{
		ID:                "",
		Name:              project.Name,
		TestRailProjectID: project.TestRailProjectID,
		CreatedAt:         time.Now().Unix(),
		UpdatedAt:         0,
	}
	oid, err := storage.Create(ctx, createProjectData)
	if err != nil {
		return Project{}, apperror.NewAppError(err, err.Error(), "")
	}
	createProjectData.ID = oid
	return createProjectData, err
}

func (s *Service) GetProjectsList(ctx context.Context, storage Storage) ([]Project, error) {
	projects, err := storage.FindAll(ctx)
	if err != nil {
		return nil, apperror.NewAppError(err, err.Error(), "")
	}
	if len(projects) == 0 {
		return []Project{}, nil
	}
	return projects, nil
}

func (s *Service) GetProjectByID(ctx context.Context, storage Storage, id string) (Project, error) {
	project, err := storage.FindOne(ctx, id)
	if err != nil {
		return Project{}, apperror.ErrNotFound
	}
	return project, nil
}

func (s *Service) UpdateProject(ctx context.Context, storage Storage, project Project) (Project, error) {
	err := project.Validate()
	if err != nil {
		return Project{}, apperror.NewAppError(err, err.Error(), "")
	}
	findProject, err := storage.FindByField(ctx, "name", project.Name)
	if err != nil {
		return Project{}, apperror.NewAppError(err, fmt.Sprintf("failed to find project with name: %s. error: %v", project.Name, err), "")
	}
	if findProject.Name != "" {
		return Project{}, apperror.NewAppError(nil, "this name is already busy", "")
	}
	findProject, err = storage.FindOne(ctx, project.ID)
	if err != nil {
		return Project{}, apperror.ErrNotFound
	}
	if project.Name != "" && findProject.Name != project.Name {
		findProject.Name = project.Name
	}
	if project.TestRailProjectID != "" && findProject.TestRailProjectID != project.TestRailProjectID {
		findProject.TestRailProjectID = project.TestRailProjectID
	}
	findProject.UpdatedAt = time.Now().Unix()
	err = storage.Update(ctx, findProject)
	if err != nil {
		return Project{}, apperror.NewAppError(err, "failed to update project", "")
	}
	return findProject, nil
}

func (s *Service) DeleteProject(ctx context.Context, storage Storage, id string) error {
	err := storage.Delete(ctx, id)
	if err != nil {
		return apperror.ErrNotFound
	}
	return nil
}
