package service

import (
	"context"
	"fmt"
	"minitrello/pkg/models/tasks"
)

func (s *Service) CreateTask(ctx context.Context, input *tasks.TaskInput) (int, error) {
	newTask := &tasks.Task{
		Name:        input.Name,
		Description: &input.Description,
		Status:      &input.Status,
		ProjectID:   input.ProjectID,
		AssignedTo:  &input.AssignedTo,
		Deadline:    &input.Deadline,
	}

	id, err := s.repo.CreateTask(ctx, newTask)
	if err != nil {
		return 0, fmt.Errorf("failed to create task: %w", err)
	}

	return id, nil
}

func (s *Service) GetTaskByID(ctx context.Context, taskID int) (*tasks.Task, error) {
	task, err := s.repo.GetTaskByID(ctx, taskID)
	if err != nil {
		return nil, fmt.Errorf("failed to get task: %w", err)
	}
	return task, nil
}

func (s *Service) UpdateTask(ctx context.Context, task *tasks.Task) error {
    err := s.repo.UpdateTask(ctx, task)
    if err != nil {
        return fmt.Errorf("failed to update task: %w", err)
    }
    return nil
}

func (s *Service) DeleteTask(ctx context.Context, taskID int) error {
    return s.repo.DeleteTask(ctx, taskID)
}

func (s *Service) ListTasksByUser(ctx context.Context, userID int) ([]*tasks.Task, error) {
    tasksList, err := s.repo.ListTasksByUser(ctx, userID)
    if err != nil {
        return nil, fmt.Errorf("failed to list tasks by user: %w", err)
    }
    return tasksList, nil
}
