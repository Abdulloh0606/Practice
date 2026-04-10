package service

import (
	"context"
	"minitrello/pkg/models"
	"minitrello/pkg/models/tasks"
	"minitrello/pkg/models/projects"
)

type IRepository interface {
	//User
	CreateUser(ctx context.Context, user *models.User) error
	GetByID(ctx context.Context, id int) (*models.User, error)
	GetByEmail(ctx context.Context, email string) (*models.User, error)
	UpdateUserName(ctx context.Context, userID int, newName string) error
	DeleteUser(ctx context.Context, id int) error
	//Project
	CreateProject(ctx context.Context, project *projects.Project) (int, error)
	GetProjectByID(ctx context.Context, id int) (*projects.Project, error)
	DeleteProject(ctx context.Context, projectID int) error

	GetUserProjectRole(ctx context.Context, projectID int, userID int) (string, error)

	AddMember(ctx context.Context, projectID int, userID int, role string) error
	DeleteMember(ctx context.Context, projectID, userID int) error
	GetUserProjects(ctx context.Context, userID int) ([]*projects.Project, error)
	//Tasks
	CreateTask(ctx context.Context, task *tasks.Task) (int, error)
	GetTaskByID(ctx context.Context, id int) (*tasks.Task, error)
	UpdateTask(ctx context.Context, task *tasks.Task) error
	DeleteTask(ctx context.Context, taskID int) error
	ListTasksByUser(ctx context.Context, userID int) ([]*tasks.Task, error)
}

type Service struct {
	repo IRepository
}

func NewService(repo IRepository) *Service {
	return &Service{
		repo: repo,
	}
}
