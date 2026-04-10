package service

import (
	"context"
	"fmt"
	"minitrello/pkg/models/projects"
)

func (s *Service) CreateProject(ctx context.Context, project *projects.ProjectInput, userID int) (int, error) {
	newProject := &projects.Project{
		Name:      project.Name,
		CreatedBy: userID,
	}
	id, err := s.repo.CreateProject(ctx, newProject)
	if err != nil {
		return 0, fmt.Errorf("failed to create project: %w", err)
	}
	err = s.repo.AddMember(ctx, id, userID, "owner")
	if err!=nil{
		return 0, fmt.Errorf("failed to add owner into project: %w", err)
	}
	return id, nil

}
func (s *Service)GetUserProjectRole(ctx context.Context, projectID int, userID int) (string, error){
	return s.repo.GetUserProjectRole(ctx, projectID, userID)
}



func (s *Service) GetProjectByID(ctx context.Context, id int) (*projects.Project, error) {
	project, err := s.repo.GetProjectByID(ctx, id)
	if err != nil{
		return nil, fmt.Errorf("failed to get project: %w", err)
	}
	return project, nil
}

func (s *Service) DeleteProject(ctx context.Context, ProjectID int) error {
	return s.repo.DeleteProject(ctx, ProjectID)
}
//добавление учатсника в проект
func (s *Service) AddMember(ctx context.Context, projectID int, userID int, role string) error {
	return s.repo.AddMember(ctx, projectID, userID, role)
}
//удаление участника из проекта
func (s *Service) DeleteMember(ctx context.Context, projectID, userID int) error {
	return s.repo.DeleteMember(ctx,projectID, userID)
}
//получение проектов где участвует пользователь
func (s *Service) GetUserProjects(ctx context.Context, userID int) ([]*projects.Project, error) {
	return s.repo.GetUserProjects(ctx, userID)
}

