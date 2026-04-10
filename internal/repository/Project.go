package repository

import (
	"context"
	"errors"
	"fmt"
	"minitrello/pkg/errs"
	"minitrello/pkg/models/projects"

	"github.com/jackc/pgx/v5"
)

// Создание проекта
func (r *Repository) CreateProject(ctx context.Context, project *projects.Project) (int, error) {
	query := `
		insert into projects(name, created_by) 
		values($1, $2)
		returning id;
	`
	var ProjectID int
	err := r.postgres.QueryRow(ctx, query, project.Name, project.CreatedBy).Scan(&ProjectID)
	if err != nil {
		return 0, err
	}
	return ProjectID, nil
}

// получение роли пользователя в проекте
func (r *Repository) GetUserProjectRole(ctx context.Context, projectID int, userID int) (string, error) {
	query := `
        select role from project_members 
        where project_id=$1 AND user_id=$2
    `
	var role string
	err := r.postgres.QueryRow(ctx, query, projectID, userID).Scan(&role)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return "", nil
		}
		return "", err
	}
	return role, nil
}

// получение проекта по ID
func (r *Repository) GetProjectByID(ctx context.Context, id int) (*projects.Project, error) {
	query := `
		select * from projects where id=$1
	`
	var project projects.Project
	err := r.postgres.QueryRow(ctx, query, id).Scan(
		&project.ID,
		&project.Name,
		&project.CreatedBy,
		&project.CreatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, errs.ErrProjectNotFound
		}
		return nil, err
	}
	return &project, nil
}

func (r *Repository) DeleteProject(ctx context.Context, projectID int) error {
	query := `delete from projects where id=$1`
	_, err := r.postgres.Exec(ctx, query, projectID)
	if err != nil {
		return err
	}
	return nil
}

func (r *Repository) AddMember(ctx context.Context, projectID int, userID int, role string) error {
	//проверяем чтобы пользователя не было в проекте
	query := `
		select 1 from project_members 
		where project_id=$1 and user_id=$2
	`
	var exist int
	err := r.postgres.QueryRow(ctx, query, projectID, userID).Scan(&exist)
	if err == nil {
		return fmt.Errorf("user already in project")
	}
	if !errors.Is(err, pgx.ErrNoRows) {
		return err
	}
	//если все нормально - добавим пользователя как участник проекта

	query = `
		insert into project_members(project_id, user_id, role)
		values($1, $2, $3)
	`
	_, err = r.postgres.Exec(ctx, query, projectID, userID, role)
	if err != nil {
		return err
	}
	return nil
}
func (r *Repository) DeleteMember(ctx context.Context, projectID, userID int) error {
	query := `
		select 1 from project_members 
		where project_id=$1 and user_id=$2
	`
	var exist int
	err := r.postgres.QueryRow(ctx, query, projectID, userID).Scan(&exist)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return errs.ErrUserNotFound
		}
		return err
	}

	query = `
		delete from project_members
		where project_id=$1 and user_id=$2
	`
	_, err = r.postgres.Exec(ctx, query, projectID, userID)
	if err != nil {
		return err
	}
	return nil
}

func (r *Repository) GetUserProjects(ctx context.Context, userID int) ([]*projects.Project, error) {
	query := `
		select p.id, p.name, p.created_by, p.created_at
		from projects p
		join project_members pm ON pm.project_id = p.idp
		where pm.user_id = $1
	`
	rows, err := r.postgres.Query(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var projectsList []*projects.Project
	for rows.Next() {
		var p projects.Project
		if err := rows.Scan(&p.ID, &p.Name, &p.CreatedBy, &p.CreatedAt); err != nil {
			return nil, err
		}
		projectsList = append(projectsList, &p)
	}
	return projectsList, nil
}
